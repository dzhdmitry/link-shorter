CREATE TABLE IF NOT EXISTS links (
    id bigserial PRIMARY KEY,
    key varchar(12) NOT NULL,
    url varchar(2000) NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS links_key_idx ON links USING hash(key);
CREATE INDEX IF NOT EXISTS links_url_idx ON links USING hash(url);
