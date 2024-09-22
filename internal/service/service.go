package service

import (
	"context"
	"go-clean-template/internal/domain"
	"go-clean-template/pkg/logger"
)

type service struct {
	lg logger.Logger
}

func NewService(lg logger.Logger) *service {
	return &service{lg: lg}
}

func (s *service) Do(ctx context.Context, req domain.ServiceRequest) error {
	s.lg.Info("Do service")
	return nil
}

func (s *service) Persist(ctx context.Context, dt string) error {
	return nil
}
