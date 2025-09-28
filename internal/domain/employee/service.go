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

	result, err := s.employeeRepo.Create(ctx, &CreateEmployeeInput{
		ProjectID:  projectID,
		Name:       req.Name,
		Email:      req.Email,
		Role:       req.Role,
		Department: req.Department,
		Status:     EmployeeStatusFromProto(req.Status),
	})
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

	return &altalunev1.CreateEmployeeResponse{
		Employee: &altalunev1.Employee{
			Id:         result.PublicID,
			Name:       result.Name,
			Email:      result.Email,
			Role:       result.Role,
			Department: result.Department,
			Status:     EmployeeStatusToProto(result.Status),
			CreatedAt:  timestamppb.New(result.CreatedAt),
			UpdatedAt:  timestamppb.New(result.UpdatedAt),
		},
		Message: "Employee created successfully",
	}, nil
}

func (s *Service) GetEmployee(ctx context.Context, req *altalunev1.GetEmployeeRequest) (*altalunev1.GetEmployeeResponse, error) {
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

	// Get employee by ID
	employee, err := s.employeeRepo.GetByID(ctx, projectID, req.EmployeeId)
	if err != nil {
		if err == ErrEmployeeNotFound {
			return nil, altalune.NewEmployeeNotFoundError(req.EmployeeId)
		}
		s.log.Error("failed to get employee",
			"error", err,
			"project_id", projectID,
			"employee_id", req.EmployeeId,
		)
		return nil, altalune.NewUnexpectedError("failed to get employee: %w", err)
	}

	return &altalunev1.GetEmployeeResponse{
		Employee: employee.ToEmployeeProto(),
	}, nil
}

func (s *Service) UpdateEmployee(ctx context.Context, req *altalunev1.UpdateEmployeeRequest) (*altalunev1.UpdateEmployeeResponse, error) {
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

	// Check if employee exists
	existingEmployee, err := s.employeeRepo.GetByID(ctx, projectID, req.EmployeeId)
	if err != nil {
		if err == ErrEmployeeNotFound {
			return nil, altalune.NewEmployeeNotFoundError(req.EmployeeId)
		}
		s.log.Error("failed to check employee existence",
			"error", err,
			"project_id", projectID,
			"employee_id", req.EmployeeId,
		)
		return nil, altalune.NewUnexpectedError("failed to check employee existence: %w", err)
	}

	// Check email uniqueness (exclude current employee)
	if req.Email != existingEmployee.Email {
		emailEmployee, err := s.employeeRepo.GetByEmail(ctx, projectID, req.Email)
		if err != nil && err != ErrEmployeeNotFound {
			s.log.Error("failed to check email uniqueness",
				"error", err,
				"project_id", projectID,
				"email", req.Email,
			)
			return nil, altalune.NewUnexpectedError("failed to check email uniqueness: %w", err)
		}
		if emailEmployee != nil && emailEmployee.ID != existingEmployee.ID {
			return nil, altalune.NewAlreadyExistsError(req.Email)
		}
	}

	// Update employee
	result, err := s.employeeRepo.Update(ctx, &UpdateEmployeeInput{
		ProjectID:  projectID,
		PublicID:   req.EmployeeId,
		Name:       req.Name,
		Email:      req.Email,
		Role:       req.Role,
		Department: req.Department,
		Status:     EmployeeStatusFromProto(req.Status),
	})

	if err != nil {
		if err == ErrEmployeeNotFound {
			return nil, altalune.NewEmployeeNotFoundError(req.EmployeeId)
		}
		if err == ErrEmployeeAlreadyExists {
			return nil, altalune.NewAlreadyExistsError(req.Email)
		}
		s.log.Error("failed to update employee",
			"error", err,
			"project_id", projectID,
			"employee_id", req.EmployeeId,
		)
		return nil, altalune.NewUnexpectedError("failed to update employee: %w", err)
	}

	return &altalunev1.UpdateEmployeeResponse{
		Employee: &altalunev1.Employee{
			Id:         result.PublicID,
			Name:       result.Name,
			Email:      result.Email,
			Role:       result.Role,
			Department: result.Department,
			Status:     EmployeeStatusToProto(result.Status),
			CreatedAt:  timestamppb.New(result.CreatedAt),
			UpdatedAt:  timestamppb.New(result.UpdatedAt),
		},
		Message: "Employee updated successfully",
	}, nil
}

func (s *Service) DeleteEmployee(ctx context.Context, req *altalunev1.DeleteEmployeeRequest) (*altalunev1.DeleteEmployeeResponse, error) {
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

	// Delete employee
	err = s.employeeRepo.Delete(ctx, &DeleteEmployeeInput{
		ProjectID: projectID,
		PublicID:  req.EmployeeId,
	})

	if err != nil {
		if err == ErrEmployeeNotFound {
			return nil, altalune.NewEmployeeNotFoundError(req.EmployeeId)
		}
		s.log.Error("failed to delete employee",
			"error", err,
			"project_id", projectID,
			"employee_id", req.EmployeeId,
		)
		return nil, altalune.NewUnexpectedError("failed to delete employee: %w", err)
	}

	return &altalunev1.DeleteEmployeeResponse{
		Message: "Employee deleted successfully",
	}, nil
}
