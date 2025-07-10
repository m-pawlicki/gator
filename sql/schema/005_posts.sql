-- +goose Up
CREATE TABLE posts (
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
title Text NOT NULL,
url Text UNIQUE NOT NULL,
description Text NOT NULL,
published_at TIMESTAMP NOT NULL,
feed_id UUID REFERENCES feeds (id) NOT NULL
);

-- +goose Down
DROP TABLE posts;