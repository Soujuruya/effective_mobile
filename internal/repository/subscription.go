package repository

import (
	"context"
	"effective_mobile/internal/models"
	"effective_mobile/pkg/logger"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	db    *pgxpool.Pool
	query squirrel.StatementBuilderType
	log   logger.Logger
}

func NewRepository(db *pgxpool.Pool, env string) *Repository {
	return &Repository{
		db:    db,
		query: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		log:   logger.NewLogger(env),
	}
}

func (r *Repository) Select(ctx context.Context, limit, offset int) []models.Subscription {
	sql, args, err := r.query.
		Select("id", "name", "price", "user_id", "start_date", "end_date").
		From("subscriptions").
		OrderBy("id").
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		r.log.Error(ctx, "Repository.Select: build query failed:", zap.Error(err))
		return []models.Subscription{}
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		r.log.Error(ctx, "Repository.Select: query failed:", zap.Error(err))
		return []models.Subscription{}
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(&s.ID, &s.Name, &s.Price, &s.UserID, &s.StartDate, &s.EndDate); err != nil {
			r.log.Error(ctx, "Repository.Select: scan failed:", zap.Error(err))
			continue
		}
		subs = append(subs, s)
	}

	return subs
}

func (r *Repository) SelectByNameAndUserID(ctx context.Context, name string, id uuid.UUID) (models.Subscription, error) {
	sql, args, err := r.query.
		Select("id", "name", "price", "user_id", "start_date", "end_date").
		From("subscriptions").
		Where(squirrel.Eq{"name": name, "user_id": id}).
		ToSql()
	if err != nil {
		r.log.Error(ctx, "Repository.SelectByNameAndUserID: builder failed", zap.Error(err))
		return models.Subscription{}, err
	}

	r.log.Debug(ctx, "Repository.SelectByNameAndUserID: executing SQL",
		zap.String("sql", sql),
		zap.Any("args", args))

	var s models.Subscription
	err = r.db.QueryRow(ctx, sql, args...).Scan(&s.ID, &s.Name, &s.Price, &s.UserID, &s.StartDate, &s.EndDate)
	if errors.Is(err, pgx.ErrNoRows) {
		r.log.Info(ctx, "Repository.SelectByNameAndUserID: query failed", zap.String("name", name), zap.String("user_id", id.String()))
		return models.Subscription{}, nil
	}
	return s, err
}

func (r *Repository) Insert(ctx context.Context, subscription *models.Subscription) error {
	sql, args, err := r.query.
		Insert("subscriptions").
		Columns("name", "price", "user_id", "start_date", "end_date").
		Values(subscription.Name, subscription.Price, subscription.UserID, subscription.StartDate, subscription.EndDate).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		r.log.Error(ctx, "Repository.Insert: builder failed", zap.Error(err))
		return err
	}

	r.log.Debug(ctx, "Repository.Insert: executing SQL",
		zap.String("sql", sql),
		zap.Any("args", args),
	)

	return r.db.QueryRow(ctx, sql, args...).Scan(&subscription.ID)
}

func (r *Repository) Update(ctx context.Context, subscription models.Subscription) error {
	sql, args, err := r.query.
		Update("subscriptions").
		Set("price", subscription.Price).
		Set("start_date", subscription.StartDate).
		Set("end_date", subscription.EndDate).
		Where(squirrel.Eq{"name": subscription.Name, "user_id": subscription.UserID}).
		ToSql()

	if err != nil {
		r.log.Error(ctx, "Repository.Update: builder failed", zap.Error(err))
		return err
	}

	r.log.Debug(ctx, "Repository.Update: executing SQL",
		zap.String("sql", sql),
		zap.Any("args", args))

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) Delete(ctx context.Context, name string, id uuid.UUID) error {
	sql, args, err := r.query.
		Delete("subscriptions").
		Where(squirrel.Eq{"name": name, "user_id": id}).
		ToSql()

	if err != nil {
		r.log.Error(ctx, "Repository.Delete: builder failed", zap.Error(err))
		return err
	}

	r.log.Debug(ctx, "Repository.Delete: executing SQL",
		zap.String("sql", sql),
		zap.Any("args", args))

	_, err = r.db.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) SumPrice(ctx context.Context, name string, id uuid.UUID, startDate string, endDate string) (int, error) {
	builder := r.query.
		Select("COALESCE(SUM(price), 0)").
		From("subscriptions").
		Where(squirrel.Eq{"user_id": id}).
		Where(squirrel.GtOrEq{"start_date": startDate})

	if endDate != "" {
		builder = builder.Where(squirrel.LtOrEq{"end_date": endDate})
	}

	if name != "" {
		builder = builder.Where(squirrel.Eq{"name": name})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		r.log.Error(ctx, "Repository.SumPrice: builder failed", zap.Error(err))
		return 0, err
	}
	r.log.Debug(ctx, "Repository.SumPrice: executing SQL",
		zap.String("sql", sql),
		zap.Any("args", args))

	var sum int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&sum)
	return sum, err
}
