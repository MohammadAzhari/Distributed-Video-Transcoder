CREATE TYPE video_status AS ENUM ('new', 'processing', 'done');

CREATE TABLE "videos" (
    "id" uuid PRIMARY KEY,
    "filename" varchar NOT NULL,
    "status" video_status NOT NULL,
    "worker_ip" varchar,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now()
);
