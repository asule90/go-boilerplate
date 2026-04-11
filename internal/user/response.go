package user

import (
	"time"

	"github.com/sule/go-boilerplate/internal/utils"
)

// UserResponse is the API response representation of a user.
type UserResponse struct {
	ID          string    `json:"id"           example:"550e8400-e29b-41d4-a716-446655440000"`
	FirebaseUID string    `json:"firebase_uid"  example:"abc123firebase"`
	DisplayName string    `json:"display_name"  example:"John Doe"`
	Email       string    `json:"email"         example:"john@example.com"`
	PhotoURL    string    `json:"photo_url"     example:"https://example.com/photo.jpg"`
	Bio         string    `json:"bio"           example:"Software developer"`
	City        string    `json:"city"          example:"Istanbul"`
	Country     string    `json:"country"       example:"Turkey"`
	CreatedAt   time.Time `json:"created_at"    example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at"    example:"2024-01-01T00:00:00Z"`
}

// UserListResponse wraps a list of users with pagination metadata.
type UserListResponse struct {
	Users      []UserResponse       `json:"users"`
	Pagination utils.PaginationMeta `json:"pagination"`
}

// ToResponse converts a User entity to UserResponse.
func ToResponse(u User) UserResponse {
	return UserResponse{
		ID:          u.ID,
		FirebaseUID: u.FirebaseUID,
		DisplayName: u.DisplayName,
		Email:       u.Email,
		PhotoURL:    u.PhotoURL,
		Bio:         u.Bio,
		City:        u.City,
		Country:     u.Country,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
