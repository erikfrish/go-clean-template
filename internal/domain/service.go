package domain

import "context"

type Service interface {
	Do(ctx context.Context, req ServiceRequest) error
	Persist(ctx context.Context, dt string) error
}
