package db

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type dbExecutor interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

// QueryBuilder wraps squirrel with convenience methods.
type QueryBuilder struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
	sq     sq.StatementBuilderType
}

// NewQueryBuilder creates a QueryBuilder with PostgreSQL dollar-sign placeholders.
func NewQueryBuilder(pool *pgxpool.Pool, logger *zap.Logger) *QueryBuilder {
	return &QueryBuilder{
		pool:   pool,
		logger: logger,
		sq:     sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (qb *QueryBuilder) getExecutor(ctx context.Context) dbExecutor {
	if tx, ok := ExtractPgxTx(ctx); ok {
		return tx
	}
	return qb.pool
}

// BaseQuery returns a SELECT builder with soft-delete filter applied.
func (qb *QueryBuilder) BaseQuery(table string, columns ...string) sq.SelectBuilder {
	return qb.sq.Select(columns...).From(table).Where(sq.Eq{table + ".deleted_at": nil})
}

// BaseQueryAll returns a SELECT builder without soft-delete filter.
func (qb *QueryBuilder) BaseQueryAll(table string, columns ...string) sq.SelectBuilder {
	return qb.sq.Select(columns...).From(table)
}

// BaseUpdate returns an UPDATE builder with updated_at = NOW().
func (qb *QueryBuilder) BaseUpdate(table string) sq.UpdateBuilder {
	return qb.sq.Update(table).Set("updated_at", sq.Expr("NOW()"))
}

// BaseInsert returns an INSERT builder for the given table.
func (qb *QueryBuilder) BaseInsert(table string) sq.InsertBuilder {
	return qb.sq.Insert(table)
}

// ApplySearch adds ILIKE conditions across multiple fields.
func (qb *QueryBuilder) ApplySearch(query sq.SelectBuilder, search string, fields []string) sq.SelectBuilder {
	if search == "" || len(fields) == 0 {
		return query
	}
	or := sq.Or{}
	for _, f := range fields {
		or = append(or, sq.ILike{f: "%" + search + "%"})
	}
	return query.Where(or)
}

// ApplySort applies an ORDER BY clause with safelist validation.
// Field may be prefixed with "-" for DESC. Defaults to "created_at DESC".
func (qb *QueryBuilder) ApplySort(query sq.SelectBuilder, field string, validSorts []string) sq.SelectBuilder {
	if field == "" {
		return query.OrderBy("created_at DESC")
	}

	dir := "ASC"
	col := field
	if len(field) > 0 && field[0] == '-' {
		dir = "DESC"
		col = field[1:]
	}

	for _, valid := range validSorts {
		if col == valid {
			return query.OrderBy(fmt.Sprintf("%s %s", col, dir))
		}
	}

	return query.OrderBy("created_at DESC")
}

// ApplyPagination applies LIMIT and OFFSET.
func (qb *QueryBuilder) ApplyPagination(query sq.SelectBuilder, page, pageSize int) sq.SelectBuilder {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	return query.Limit(uint64(pageSize)).Offset(uint64(offset))
}

// ExecuteQuery runs a SELECT query and returns multiple rows.
func (qb *QueryBuilder) ExecuteQuery(ctx context.Context, query sq.SelectBuilder) (pgx.Rows, error) {
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query builder ToSql: %w", err)
	}
	qb.logQuery(sqlStr, args)
	return qb.getExecutor(ctx).Query(ctx, sqlStr, args...)
}

// ExecuteQueryRow runs a SELECT query and returns a single row.
func (qb *QueryBuilder) ExecuteQueryRow(ctx context.Context, query sq.SelectBuilder) pgx.Row {
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return &errorRow{err: err}
	}
	qb.logQuery(sqlStr, args)
	return qb.getExecutor(ctx).QueryRow(ctx, sqlStr, args...)
}

// ExecuteCount runs a COUNT(*) query and returns the count.
func (qb *QueryBuilder) ExecuteCount(ctx context.Context, query sq.SelectBuilder) (int64, error) {
	countQuery := qb.sq.Select("COUNT(*)").FromSelect(query, "sub")
	sqlStr, args, err := countQuery.ToSql()
	if err != nil {
		return 0, fmt.Errorf("count query ToSql: %w", err)
	}
	qb.logQuery(sqlStr, args)
	var count int64
	err = qb.getExecutor(ctx).QueryRow(ctx, sqlStr, args...).Scan(&count)
	return count, err
}

// ExecuteInsert runs an INSERT and scans the result into dest.
func (qb *QueryBuilder) ExecuteInsert(ctx context.Context, query sq.InsertBuilder, dest ...interface{}) error {
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("insert ToSql: %w", err)
	}
	qb.logQuery(sqlStr, args)
	return qb.getExecutor(ctx).QueryRow(ctx, sqlStr, args...).Scan(dest...)
}

// ExecuteUpdate runs an UPDATE and scans the result into dest.
func (qb *QueryBuilder) ExecuteUpdate(ctx context.Context, query sq.UpdateBuilder, dest ...interface{}) error {
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("update ToSql: %w", err)
	}
	qb.logQuery(sqlStr, args)
	if len(dest) == 0 {
		_, err = qb.getExecutor(ctx).Exec(ctx, sqlStr, args...)
		return err
	}
	return qb.getExecutor(ctx).QueryRow(ctx, sqlStr, args...).Scan(dest...)
}

// ExecuteDelete runs a soft DELETE (UPDATE with deleted_at) statement.
func (qb *QueryBuilder) ExecuteDelete(ctx context.Context, query sq.UpdateBuilder) error {
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("delete ToSql: %w", err)
	}
	qb.logQuery(sqlStr, args)
	_, err = qb.getExecutor(ctx).Exec(ctx, sqlStr, args...)
	return err
}

func (qb *QueryBuilder) logQuery(sql string, args []interface{}) {
	if qb.logger != nil {
		qb.logger.Debug("executing query", zap.String("sql", sql), zap.Any("args", args))
	}
}

type errorRow struct {
	err error
}

func (e *errorRow) Scan(_ ...interface{}) error {
	return e.err
}
