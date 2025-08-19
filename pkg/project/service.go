// pkg/project/service.go
package project

import (
	"context"
	"time"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/query"
)

type Service struct {
	altalunev1.UnimplementedProjectServiceServer
	validator   protovalidate.Validator
	log         altalune.Logger
	projectRepo Repositor
}

func NewService(v protovalidate.Validator, log altalune.Logger, projectRepo Repositor) *Service {
	return &Service{
		validator:   v,
		log:         log,
		projectRepo: projectRepo,
	}
}

func (s *Service) QueryProjects(ctx context.Context, req *altalunev1.QueryProjectsRequest) (*altalunev1.QueryProjectsResponse, error) {
	// Add artificial delay to match frontend pattern
	time.Sleep(2000 * time.Millisecond)

	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Set default query if not provided
	if req.Query == nil {
		req.Query = &altalunev1.QueryRequest{
			Pagination: &altalunev1.Pagination{
				Page:     1,
				PageSize: 10,
			},
		}
	}

	// Convert proto request to domain query params
	queryParams := query.DefaultQueryParams(req.Query)

	// Query projects from repository
	result, err := s.projectRepo.Query(ctx, queryParams)
	if err != nil {
		s.log.Error("failed to query projects",
			"error", err,
			"keyword", queryParams.Keyword,
		)
		return nil, altalune.NewUnexpectedError("failed to query projects: %w", err)
	}

	// Convert domain result to proto response
	if result == nil {
		return &altalunev1.QueryProjectsResponse{
			Data: []*altalunev1.Project{},
			Meta: &altalunev1.QueryMetaResponse{
				RowCount:  0,
				PageCount: 0,
				Filters:   make(map[string]*altalunev1.FilterValues),
			},
		}, nil
	}

	return &altalunev1.QueryProjectsResponse{
		Data: mapProjectsToProto(result.Data),
		Meta: &altalunev1.QueryMetaResponse{
			RowCount:  result.TotalRows,
			PageCount: result.TotalPages,
			Filters:   mapFiltersToProto(result.Filters),
		},
	}, nil
}
