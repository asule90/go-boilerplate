package user

// CreateUserRequest is used to create or upsert a user via Firebase.
type CreateUserRequest struct {
	FirebaseUID string `json:"firebase_uid" validate:"required"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=255"`
	Email       string `json:"email"        validate:"required,email"`
	PhotoURL    string `json:"photo_url"`
}

// UpdateUserRequest is used to update a user's profile.
type UpdateUserRequest struct {
	DisplayName string `json:"display_name" validate:"omitempty,min=1,max=255"`
	Bio         string `json:"bio"          validate:"omitempty,max=500"`
	City        string `json:"city"         validate:"omitempty,max=100"`
	Country     string `json:"country"      validate:"omitempty,max=100"`
	PhotoURL    string `json:"photo_url"    validate:"omitempty,url"`
}

// ListUsersRequest holds pagination, search, and sort parameters for listing users.
type ListUsersRequest struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Search   string `query:"search"`
	Sort     string `query:"sort"`
}
