package user

import (
	"context"
	"strings"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"
	altalunev1 "github.com/hrz8/altalune/gen/altalune/v1"
	"github.com/hrz8/altalune/internal/shared/query"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	altalunev1.UnimplementedUserServiceServer
	validator protovalidate.Validator
	log       altalune.Logger
	userRepo  Repository
}

func NewService(v protovalidate.Validator, log altalune.Logger, userRepo Repository) *Service {
	return &Service{
		validator: v,
		log:       log,
		userRepo:  userRepo,
	}
}

func (s *Service) QueryUsers(ctx context.Context, req *altalunev1.QueryUsersRequest) (*altalunev1.QueryUsersResponse, error) {
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

	// Query users from repository
	result, err := s.userRepo.Query(ctx, queryParams)
	if err != nil {
		s.log.Error("failed to query users",
			"error", err,
			"keyword", queryParams.Keyword,
		)
		return nil, altalune.NewUnexpectedError("failed to query users: %w", err)
	}

	// Convert domain result to proto response
	if result == nil {
		return &altalunev1.QueryUsersResponse{
			Data: []*altalunev1.User{},
			Meta: &altalunev1.QueryMetaResponse{
				RowCount:  0,
				PageCount: 0,
				Filters:   make(map[string]*altalunev1.FilterValues),
			},
		}, nil
	}

	return &altalunev1.QueryUsersResponse{
		Data: mapUsersToProto(result.Data),
		Meta: &altalunev1.QueryMetaResponse{
			RowCount:  result.TotalRows,
			PageCount: result.TotalPages,
			Filters:   mapFiltersToProto(result.Filters),
		},
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, req *altalunev1.CreateUserRequest) (*altalunev1.CreateUserResponse, error) {
	// Validate request
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Lowercase email for consistency
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check if user with same email already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && err != ErrUserNotFound {
		s.log.Error("failed to check existing user",
			"error", err,
			"email", email,
		)
		return nil, altalune.NewUnexpectedError("failed to check existing user: %w", err)
	}

	if existingUser != nil {
		return nil, altalune.NewUserAlreadyExistsError(email)
	}

	result, err := s.userRepo.Create(ctx, &CreateUserInput{
		Email:     email,
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
	})
	if err != nil {
		if err == ErrUserAlreadyExists {
			return nil, altalune.NewUserAlreadyExistsError(email)
		}
		s.log.Error("failed to create user",
			"error", err,
			"email", email,
		)
		return nil, altalune.NewUnexpectedError("failed to create user: %w", err)
	}

	return &altalunev1.CreateUserResponse{
		User: &altalunev1.User{
			Id:        result.PublicID,
			Email:     result.Email,
			FirstName: result.FirstName,
			LastName:  result.LastName,
			IsActive:  result.IsActive,
			CreatedAt: timestamppb.New(result.CreatedAt),
			UpdatedAt: timestamppb.New(result.UpdatedAt),
		},
		Message: "User created successfully",
	}, nil
}

func (s *Service) GetUser(ctx context.Context, req *altalunev1.GetUserRequest) (*altalunev1.GetUserResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	user, err := s.userRepo.GetByID(ctx, req.Id)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, altalune.NewUserNotFoundError(req.Id)
		}
		s.log.Error("failed to get user", "error", err, "user_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to get user", err)
	}

	return &altalunev1.GetUserResponse{
		User: user.ToUserProto(),
	}, nil
}

func (s *Service) UpdateUser(ctx context.Context, req *altalunev1.UpdateUserRequest) (*altalunev1.UpdateUserResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	// Get internal ID
	internalID, err := s.userRepo.GetIDByPublicID(ctx, req.Id)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, altalune.NewUserNotFoundError(req.Id)
		}
		return nil, altalune.NewUnexpectedError("failed to resolve user ID", err)
	}

	// Lowercase email for consistency
	email := strings.ToLower(strings.TrimSpace(req.Email))

	input := &UpdateUserInput{
		ID:        internalID,
		PublicID:  req.Id,
		Email:     email,
		FirstName: strings.TrimSpace(req.FirstName),
		LastName:  strings.TrimSpace(req.LastName),
	}

	result, err := s.userRepo.Update(ctx, input)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, altalune.NewUserNotFoundError(req.Id)
		}
		if err == ErrUserAlreadyExists {
			return nil, altalune.NewUserAlreadyExistsError(email)
		}
		s.log.Error("failed to update user", "error", err, "user_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to update user", err)
	}

	s.log.Info("user updated successfully", "user_id", req.Id, "email", email)

	return &altalunev1.UpdateUserResponse{
		User:    result.ToUser().ToUserProto(),
		Message: "User updated successfully",
	}, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *altalunev1.DeleteUserRequest) (*altalunev1.DeleteUserResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	err := s.userRepo.Delete(ctx, req.Id)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, altalune.NewUserNotFoundError(req.Id)
		}
		s.log.Error("failed to delete user", "error", err, "user_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to delete user", err)
	}

	s.log.Info("user deleted successfully", "user_id", req.Id)

	return &altalunev1.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}

func (s *Service) ActivateUser(ctx context.Context, req *altalunev1.ActivateUserRequest) (*altalunev1.ActivateUserResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	user, err := s.userRepo.Activate(ctx, req.Id)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, altalune.NewUserNotFoundError(req.Id)
		}
		if err == ErrUserAlreadyActive {
			return nil, altalune.NewUserAlreadyActiveError(req.Id)
		}
		s.log.Error("failed to activate user", "error", err, "user_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to activate user", err)
	}

	s.log.Info("user activated successfully", "user_id", req.Id)

	return &altalunev1.ActivateUserResponse{
		User:    user.ToUserProto(),
		Message: "User activated successfully",
	}, nil
}

func (s *Service) DeactivateUser(ctx context.Context, req *altalunev1.DeactivateUserRequest) (*altalunev1.DeactivateUserResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}

	user, err := s.userRepo.Deactivate(ctx, req.Id)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, altalune.NewUserNotFoundError(req.Id)
		}
		if err == ErrUserAlreadyInactive {
			return nil, altalune.NewUserAlreadyInactiveError(req.Id)
		}
		s.log.Error("failed to deactivate user", "error", err, "user_id", req.Id)
		return nil, altalune.NewUnexpectedError("failed to deactivate user", err)
	}

	s.log.Info("user deactivated successfully", "user_id", req.Id)

	return &altalunev1.DeactivateUserResponse{
		User:    user.ToUserProto(),
		Message: "User deactivated successfully",
	}, nil
}
