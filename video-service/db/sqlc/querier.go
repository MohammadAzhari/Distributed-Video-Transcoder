// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error)
	GetVideo(ctx context.Context, id uuid.UUID) (Video, error)
	PublishVideo(ctx context.Context, arg PublishVideoParams) (Video, error)
}

var _ Querier = (*Queries)(nil)
