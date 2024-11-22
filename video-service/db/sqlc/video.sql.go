// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: video.sql

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createVideo = `-- name: CreateVideo :one
INSERT INTO videos (
  id, filename, status
) VALUES (
  $1, $2, 'new'
)
RETURNING id, filename, status, worker_ip, created_at, updated_at, scales
`

type CreateVideoParams struct {
	ID       uuid.UUID `json:"id"`
	Filename string    `json:"filename"`
}

func (q *Queries) CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error) {
	row := q.db.QueryRow(ctx, createVideo, arg.ID, arg.Filename)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Filename,
		&i.Status,
		&i.WorkerIp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Scales,
	)
	return i, err
}

const getVideo = `-- name: GetVideo :one
SELECT id, filename, status, worker_ip, created_at, updated_at, scales FROM videos
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetVideo(ctx context.Context, id uuid.UUID) (Video, error) {
	row := q.db.QueryRow(ctx, getVideo, id)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Filename,
		&i.Status,
		&i.WorkerIp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Scales,
	)
	return i, err
}

const publishVideo = `-- name: PublishVideo :one
UPDATE videos
SET status = 'done', worker_ip = $2, scales = $3
WHERE id = $1
RETURNING id, filename, status, worker_ip, created_at, updated_at, scales
`

type PublishVideoParams struct {
	ID       uuid.UUID   `json:"id"`
	WorkerIp pgtype.Text `json:"worker_ip"`
	Scales   []string    `json:"scales"`
}

func (q *Queries) PublishVideo(ctx context.Context, arg PublishVideoParams) (Video, error) {
	row := q.db.QueryRow(ctx, publishVideo, arg.ID, arg.WorkerIp, arg.Scales)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Filename,
		&i.Status,
		&i.WorkerIp,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Scales,
	)
	return i, err
}
