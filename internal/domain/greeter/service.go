package greeter

import (
	"context"
	"math"
	"time"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"
)

type Service struct {
	greeterv1.UnimplementedGreeterServiceServer
	validator   protovalidate.Validator
	log         altalune.Logger
	greeterRepo Repositor
}

func NewService(v protovalidate.Validator, log altalune.Logger, greeterRepo Repositor) *Service {
	return &Service{
		validator:   v,
		log:         log,
		greeterRepo: greeterRepo,
	}
}

func (s *Service) SayHello(ctx context.Context, req *greeterv1.SayHelloRequest) (*greeterv1.SayHelloResponse, error) {
	time.Sleep(700 * time.Millisecond) // Simulate some processing delay
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	allowedNameMap := getAllowedNameMap()
	if _, ok := allowedNameMap[req.Name]; !ok {
		return nil, altalune.NewGreetingUnrecognize(req.Name)
	}

	msg := s.greeterRepo.GetGreeterTemplate(req.Name)
	response := &greeterv1.SayHelloResponse{
		Message: msg,
	}
	return response, nil
}

func getAllowedNameMap() map[string]bool {
	m := make(map[string]bool, len(allowedNames))
	for _, name := range allowedNames {
		m[name] = true
	}
	return m
}

func (s *Service) GetAllowedNames(ctx context.Context, req *greeterv1.GetAllowedNamesRequest) (*greeterv1.GetAllowedNamesResponse, error) {
	time.Sleep(700 * time.Millisecond) // Simulate some processing delay

	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Get paginated names and total count
	names, total := s.greeterRepo.GetAllowedNamesWithTotal(req.Page, req.Limit)

	// Calculate pagination metadata
	totalPages := int32(math.Ceil(float64(total) / float64(req.Limit)))
	hasNext := req.Page < totalPages
	hasPrev := req.Page > 1

	meta := &greeterv1.PaginationMeta{
		Total:      total,
		Page:       req.Page,
		Limit:      req.Limit,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	return &greeterv1.GetAllowedNamesResponse{
		Names: names,
		Meta:  meta,
	}, nil
}
