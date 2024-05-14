-- +goose Up
CREATE TABLE follows (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL,
    CONSTRAINT fk_users FOREIGN KEY (user_id) REFERENCES users(id),
    feed_id UUID NOT NULL,
    CONSTRAINT fk_feeds FOREIGN KEY (feed_id) REFERENCES feeds(id)
);

-- +goose Down
DROP TABLE feeds;