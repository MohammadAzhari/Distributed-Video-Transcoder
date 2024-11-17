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
  filename, status
) VALUES (
  $1, 'new'
)
RETURNING id, filename, status, worker_ip, created_at, updated_at
`

func (q *Queries) CreateVideo(ctx context.Context, filename string) (Video, error) {
	row := q.db.QueryRow(ctx, createVideo, filename)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Filename,
		&i.Status,
		&i.WorkerIp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getVideo = `-- name: GetVideo :one
SELECT id, filename, status, worker_ip, created_at, updated_at FROM videos
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
	)
	return i, err
}

const publishVideo = `-- name: PublishVideo :one
UPDATE videos
SET status = 'done', worker_ip = $2
WHERE id = $1
RETURNING id, filename, status, worker_ip, created_at, updated_at
`

type PublishVideoParams struct {
	ID       uuid.UUID   `json:"id"`
	WorkerIp pgtype.Text `json:"worker_ip"`
}

func (q *Queries) PublishVideo(ctx context.Context, arg PublishVideoParams) (Video, error) {
	row := q.db.QueryRow(ctx, publishVideo, arg.ID, arg.WorkerIp)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Filename,
		&i.Status,
		&i.WorkerIp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
