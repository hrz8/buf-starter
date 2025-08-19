package employee

import (
	"context"
	"time"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/query"
	project_domain "github.com/hrz8/altalune/pkg/project"
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
