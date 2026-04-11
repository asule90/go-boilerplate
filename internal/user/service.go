package user

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/sule/go-boilerplate/internal/utils"
	"github.com/sule/go-boilerplate/pkg/errr"
)

// Service defines the business logic interface for users.
type Service interface {
	UpsertByFirebaseUID(ctx context.Context, req CreateUserRequest) (UserResponse, error)
	GetByID(ctx context.Context, id string) (UserResponse, error)
	List(ctx context.Context, req ListUsersRequest) (UserListResponse, error)
	UpdateProfile(ctx context.Context, id string, req UpdateUserRequest) (UserResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

// NewService creates a new user service.
func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) UpsertByFirebaseUID(ctx context.Context, req CreateUserRequest) (UserResponse, error) {
	existing, err := s.repo.GetByFirebaseUID(ctx, req.FirebaseUID)
	if err != nil && !errors.Is(err, errr.ErrNoRows) {
		return UserResponse{}, fmt.Errorf("user.svc.UpsertByFirebaseUID: %w", err)
	}

	if err == nil {
		existing.DisplayName = req.DisplayName
		existing.PhotoURL = req.PhotoURL
		updated, updateErr := s.repo.Update(ctx, existing)
		if updateErr != nil {
			return UserResponse{}, fmt.Errorf("user.svc.UpsertByFirebaseUID update: %w", updateErr)
		}
		return ToResponse(updated), nil
	}

	newUser := User{
		FirebaseUID: req.FirebaseUID,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		PhotoURL:    req.PhotoURL,
	}
	created, createErr := s.repo.Create(ctx, newUser)
	if createErr != nil {
		return UserResponse{}, fmt.Errorf("user.svc.UpsertByFirebaseUID create: %w", createErr)
	}
	return ToResponse(created), nil
}

func (s *service) GetByID(ctx context.Context, id string) (UserResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return UserResponse{}, fmt.Errorf("user.svc.GetByID: %w", err)
	}
	return ToResponse(u), nil
}

func (s *service) List(ctx context.Context, req ListUsersRequest) (UserListResponse, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	users, total, err := s.repo.List(ctx, req)
	if err != nil {
		return UserListResponse{}, fmt.Errorf("user.svc.List: %w", err)
	}

	responses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, ToResponse(u))
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))
	if totalPages < 1 {
		totalPages = 1
	}

	return UserListResponse{
		Users: responses,
		Pagination: utils.PaginationMeta{
			CurrentPage:  req.Page,
			PageSize:     req.PageSize,
			TotalPages:   totalPages,
			TotalRecords: total,
		},
	}, nil
}

func (s *service) UpdateProfile(ctx context.Context, id string, req UpdateUserRequest) (UserResponse, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return UserResponse{}, fmt.Errorf("user.svc.UpdateProfile: %w", err)
	}

	if req.DisplayName != "" {
		existing.DisplayName = req.DisplayName
	}
	if req.Bio != "" {
		existing.Bio = req.Bio
	}
	if req.City != "" {
		existing.City = req.City
	}
	if req.Country != "" {
		existing.Country = req.Country
	}
	if req.PhotoURL != "" {
		existing.PhotoURL = req.PhotoURL
	}

	updated, err := s.repo.Update(ctx, existing)
	if err != nil {
		return UserResponse{}, fmt.Errorf("user.svc.UpdateProfile update: %w", err)
	}
	return ToResponse(updated), nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("user.svc.Delete: %w", err)
	}
	return nil
}
