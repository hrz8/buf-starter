package employee

import (
	"context"
	"time"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	project_domain "github.com/hrz8/altalune/internal/domain/project"
	"github.com/hrz8/altalune/internal/shared/query"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	altalunev1.UnimplementedEmployeeServiceServer
	validator    protovalidate.Validator
	log          altalune.Logger
	projectRepo  project_domain.Repositor
	employeeRepo Repositor
}

func NewService(v protovalidate.Validator, log altalune.Logger, projectRepo project_domain.Repositor, employeeRepo Repositor) *Service {
	return &Service{
		validator:    v,
		log:          log,
		projectRepo:  projectRepo,
		employeeRepo: employeeRepo,
	}
}

func (s *Service) QueryEmployees(ctx context.Context, req *altalunev1.QueryEmployeesRequest) (*altalunev1.QueryEmployeesResponse, error) {
	// Add artificial delay to match TypeScript version
	time.Sleep(600 * time.Millisecond)

	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Validate that query is provided
	if req.Query == nil {
		return nil, altalune.NewInvalidPayloadError("query is required")
	}

	// Extract project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// Convert proto request to domain query params
	queryParams := query.DefaultQueryParams(req.Query)

	// Query employees from repository
	result, err := s.employeeRepo.Query(ctx, projectID, queryParams)
	if err != nil {
		s.log.Error("failed to query employees",
			"error", err,
			"project_id", projectID,
			"keyword", queryParams.Keyword,
		)
		return nil, altalune.NewUnexpectedError("failed to query employees: %w", err)
	}

	// Convert domain result to proto response
	if result == nil {
		return &altalunev1.QueryEmployeesResponse{
			Data: []*altalunev1.Employee{},
			Meta: &altalunev1.QueryMetaResponse{
				RowCount:  0,
				PageCount: 0,
				Filters:   make(map[string]*altalunev1.FilterValues),
			},
		}, nil
	}

	return &altalunev1.QueryEmployeesResponse{
		Data: mapEmployeesToProto(result.Data),
		Meta: &altalunev1.QueryMetaResponse{
			RowCount:  result.TotalRows,
			PageCount: result.TotalPages,
			Filters:   mapFiltersToProto(result.Filters),
		},
	}, nil
}

func (s *Service) CreateEmployee(ctx context.Context, req *altalunev1.CreateEmployeeRequest) (*altalunev1.CreateEmployeeResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Extract and validate project ID
	projectID, err := s.projectRepo.GetIDByPublicID(ctx, req.ProjectId)
	if err != nil {
		if err == project_domain.ErrProjectNotFound {
			return nil, altalune.NewProjectNotFound(req.ProjectId)
		}
		return nil, altalune.NewInvalidPayloadError("invalid project_id")
	}

	// Check if employee with same email already exists
	existingEmployee, err := s.employeeRepo.GetByEmail(ctx, projectID, req.Email)
	if err != nil && err != ErrEmployeeNotFound {
		s.log.Error("failed to check existing employee",
			"error", err,
			"project_id", projectID,
			"email", req.Email,
		)
		return nil, altalune.NewUnexpectedError("failed to check existing employee: %w", err)
	}

	if existingEmployee != nil {
		return nil, altalune.NewAlreadyExistsError(req.Email)
	}

	// Map proto status to domain status
	var domainStatus EmployeeStatus
	switch req.Status {
	case altalunev1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE:
		domainStatus = EmployeeStatusActive
	case altalunev1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE:
		domainStatus = EmployeeStatusInactive
	default:
		domainStatus = EmployeeStatusActive
	}

	// Create the employee
	input := &CreateEmployeeInput{
		ProjectID:  projectID,
		Name:       req.Name,
		Email:      req.Email,
		Role:       req.Role,
		Department: req.Department,
		Status:     domainStatus,
	}

	result, err := s.employeeRepo.Create(ctx, input)
	if err != nil {
		if err == ErrEmployeeAlreadyExists {
			return nil, altalune.NewAlreadyExistsError(req.Email)
		}
		s.log.Error("failed to create employee",
			"error", err,
			"project_id", projectID,
			"name", req.Name,
			"email", req.Email,
		)
		return nil, altalune.NewUnexpectedError("failed to create employee: %w", err)
	}

	// Map domain result to proto response
	protoStatus := altalunev1.EmployeeStatus_EMPLOYEE_STATUS_UNSPECIFIED
	switch result.Status {
	case EmployeeStatusActive:
		protoStatus = altalunev1.EmployeeStatus_EMPLOYEE_STATUS_ACTIVE
	case EmployeeStatusInactive:
		protoStatus = altalunev1.EmployeeStatus_EMPLOYEE_STATUS_INACTIVE
	}

	return &altalunev1.CreateEmployeeResponse{
		Employee: &altalunev1.Employee{
			Id:         result.PublicID,
			Name:       result.Name,
			Email:      result.Email,
			Role:       result.Role,
			Department: result.Department,
			Status:     protoStatus,
			CreatedAt:  timestamppb.New(result.CreatedAt),
			UpdatedAt:  timestamppb.New(result.UpdatedAt),
		},
		Message: "Employee created successfully",
	}, nil
}
