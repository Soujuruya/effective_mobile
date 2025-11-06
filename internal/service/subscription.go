package service

import (
	"context"
	"effective_mobile/internal/models"
	"effective_mobile/pkg/logger"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionRepository interface {
	Select(ctx context.Context, limit, offset int) []models.Subscription
	SelectByNameAndUserID(ctx context.Context, name string, id uuid.UUID) (models.Subscription, error)
	Insert(ctx context.Context, subscription *models.Subscription) error
	Update(ctx context.Context, subscription models.Subscription) error
	Delete(ctx context.Context, name string, id uuid.UUID) error
	SumPrice(ctx context.Context, name string, id uuid.UUID, startDate string, endDate string) (int, error)
}

type SubscriptionService struct {
	repo SubscriptionRepository
	log  logger.Logger
}

func NewSubscriptionService(repository SubscriptionRepository, env string) *SubscriptionService {
	return &SubscriptionService{
		repo: repository,
		log:  logger.NewLogger(env),
	}
}

func (s *SubscriptionService) Select(ctx context.Context, limit, offset int) []models.Subscription {
	s.log.Debug(ctx, "Service.Select called",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)

	subs := s.repo.Select(ctx, limit, offset)

	s.log.Debug(ctx, "Service.Select result",
		zap.Int("subscriptions_count", len(subs)),
	)
	return subs
}

func (s *SubscriptionService) SelectByNameAndUserID(ctx context.Context, name string, id uuid.UUID) (models.Subscription, error) {
	s.log.Debug(ctx, "Service.SelectByNameAndUserID called",
		zap.String("name", name),
		zap.String("user_id", id.String()),
	)

	sub, err := s.repo.SelectByNameAndUserID(ctx, name, id)
	if err != nil {
		s.log.Error(ctx, "Service.SelectByNameAndUserID error", zap.Error(err))
	} else {
		s.log.Debug(ctx, "Service.SelectByNameAndUserID result", zap.Any("subscription", sub))
	}

	return sub, err
}

func (s *SubscriptionService) Insert(ctx context.Context, subscription *models.Subscription) error {
	s.log.Debug(ctx, "Service.Insert called", zap.Any("subscription", subscription))
	err := s.repo.Insert(ctx, subscription)
	if err != nil {
		s.log.Error(ctx, "Service.Insert error", zap.Error(err))
	} else {
		s.log.Debug(ctx, "Service.Insert successful", zap.Int("subscription_id", subscription.ID))
	}
	return err
}

func (s *SubscriptionService) Update(ctx context.Context, subscription models.Subscription) error {
	s.log.Debug(ctx, "Service.Update called", zap.Any("subscription", subscription))
	err := s.repo.Update(ctx, subscription)
	if err != nil {
		s.log.Error(ctx, "Service.Update error", zap.Error(err))
	} else {
		s.log.Debug(ctx, "Service.Update successful")
	}
	return err
}

func (s *SubscriptionService) Delete(ctx context.Context, name string, id uuid.UUID) error {
	s.log.Debug(ctx, "Service.Delete called",
		zap.String("name", name),
		zap.String("user_id", id.String()),
	)

	err := s.repo.Delete(ctx, name, id)
	if err != nil {
		s.log.Error(ctx, "Service.Delete error", zap.Error(err))
	} else {
		s.log.Debug(ctx, "Service.Delete successful")
	}
	return err
}

func (s *SubscriptionService) SumPrice(ctx context.Context, name string, id uuid.UUID, startDate string, endDate string) (int, error) {
	s.log.Debug(ctx, "Service.SumPrice called",
		zap.String("name", name),
		zap.String("user_id", id.String()),
		zap.String("start_date", startDate),
		zap.String("end_date", endDate),
	)

	sum, err := s.repo.SumPrice(ctx, name, id, startDate, endDate)
	if err != nil {
		s.log.Error(ctx, "Service.SumPrice error", zap.Error(err))
	} else {
		s.log.Debug(ctx, "Service.SumPrice result", zap.Int("sum", sum))
	}

	return sum, err
}
