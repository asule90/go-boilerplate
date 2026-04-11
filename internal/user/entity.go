package user

import "time"

// User is the database entity for the users table.
type User struct {
	ID          string     `db:"id"           json:"id"`
	FirebaseUID string     `db:"firebase_uid"  json:"firebase_uid"`
	DisplayName string     `db:"display_name"  json:"display_name"`
	Email       string     `db:"email"         json:"email"`
	PhotoURL    string     `db:"photo_url"     json:"photo_url"`
	Bio         string     `db:"bio"           json:"bio"`
	City        string     `db:"city"          json:"city"`
	Country     string     `db:"country"       json:"country"`
	CreatedAt   time.Time  `db:"created_at"    json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"    json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"    json:"deleted_at,omitempty"`
}
