package user

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/sule/go-boilerplate/internal/db"
	"github.com/sule/go-boilerplate/pkg/errr"
)

const tableName = "users"

var userColumns = []string{
	"id", "firebase_uid", "display_name", "email",
	"photo_url", "bio", "city", "country",
	"created_at", "updated_at", "deleted_at",
}

// Repository defines the data access interface for users.
type Repository interface {
	Create(ctx context.Context, u User) (User, error)
	GetByID(ctx context.Context, id string) (User, error)
	GetByFirebaseUID(ctx context.Context, uid string) (User, error)
	List(ctx context.Context, req ListUsersRequest) ([]User, int64, error)
	Update(ctx context.Context, u User) (User, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	qb *db.QueryBuilder
}

// NewRepository creates a new user repository.
func NewRepository(qb *db.QueryBuilder) Repository {
	return &repository{qb: qb}
}

func (r *repository) Create(ctx context.Context, u User) (User, error) {
	query := r.qb.BaseInsert(tableName).
		Columns("firebase_uid", "display_name", "email", "photo_url").
		Values(u.FirebaseUID, u.DisplayName, u.Email, u.PhotoURL).
		Suffix("RETURNING *")

	var result User
	err := r.qb.ExecuteInsert(ctx, query,
		&result.ID, &result.FirebaseUID, &result.DisplayName, &result.Email,
		&result.PhotoURL, &result.Bio, &result.City, &result.Country,
		&result.CreatedAt, &result.UpdatedAt, &result.DeletedAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("user.repo.Create: %w", errr.ParseDBError(err))
	}
	return result, nil
}

func (r *repository) GetByID(ctx context.Context, id string) (User, error) {
	query := r.qb.BaseQuery(tableName, userColumns...).
		Where(sq.Eq{tableName + ".id": id})

	var u User
	err := r.qb.ExecuteQueryRow(ctx, query).Scan(
		&u.ID, &u.FirebaseUID, &u.DisplayName, &u.Email,
		&u.PhotoURL, &u.Bio, &u.City, &u.Country,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("user.repo.GetByID: %w", errr.ParseDBError(err))
	}
	return u, nil
}

func (r *repository) GetByFirebaseUID(ctx context.Context, uid string) (User, error) {
	query := r.qb.BaseQuery(tableName, userColumns...).
		Where(sq.Eq{tableName + ".firebase_uid": uid})

	var u User
	err := r.qb.ExecuteQueryRow(ctx, query).Scan(
		&u.ID, &u.FirebaseUID, &u.DisplayName, &u.Email,
		&u.PhotoURL, &u.Bio, &u.City, &u.Country,
		&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("user.repo.GetByFirebaseUID: %w", errr.ParseDBError(err))
	}
	return u, nil
}

func (r *repository) List(ctx context.Context, req ListUsersRequest) ([]User, int64, error) {
	validSorts := []string{"display_name", "email", "created_at", "updated_at"}
	prefixedCols := make([]string, len(userColumns))
	for i, col := range userColumns {
		prefixedCols[i] = tableName + "." + col
	}

	baseQuery := r.qb.BaseQuery(tableName, prefixedCols...)
	baseQuery = r.qb.ApplySearch(baseQuery, req.Search, []string{
		tableName + ".display_name",
		tableName + ".email",
	})
	baseQuery = r.qb.ApplySort(baseQuery, req.Sort, validSorts)

	total, err := r.qb.ExecuteCount(ctx, baseQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("user.repo.List count: %w", err)
	}

	paginatedQuery := r.qb.ApplyPagination(baseQuery, req.Page, req.PageSize)
	rows, err := r.qb.ExecuteQuery(ctx, paginatedQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("user.repo.List: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID, &u.FirebaseUID, &u.DisplayName, &u.Email,
			&u.PhotoURL, &u.Bio, &u.City, &u.Country,
			&u.CreatedAt, &u.UpdatedAt, &u.DeletedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("user.repo.List scan: %w", err)
		}
		users = append(users, u)
	}

	return users, total, nil
}

func (r *repository) Update(ctx context.Context, u User) (User, error) {
	query := r.qb.BaseUpdateWithTimestamp(tableName).
		Set("display_name", u.DisplayName).
		Set("bio", u.Bio).
		Set("city", u.City).
		Set("country", u.Country).
		Set("photo_url", u.PhotoURL).
		Where(sq.Eq{"id": u.ID}).
		Suffix("RETURNING *")

	var result User
	err := r.qb.ExecuteUpdate(ctx, query,
		&result.ID, &result.FirebaseUID, &result.DisplayName, &result.Email,
		&result.PhotoURL, &result.Bio, &result.City, &result.Country,
		&result.CreatedAt, &result.UpdatedAt, &result.DeletedAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("user.repo.Update: %w", errr.ParseDBError(err))
	}
	return result, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	query := r.qb.BaseUpdateWithTimestamp(tableName).
		Set("deleted_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": id})

	if err := r.qb.ExecuteDelete(ctx, query); err != nil {
		return fmt.Errorf("user.repo.Delete: %w", errr.ParseDBError(err))
	}
	return nil
}
