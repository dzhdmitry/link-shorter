CREATE TABLE IF NOT EXISTS links (
    id bigserial PRIMARY KEY,
    key varchar(12) NOT NULL,
    url varchar(2000) NOT NULL
);

CREATE INDEX IF NOT EXISTS links_key_idx ON links USING hash(key);
