CREATE TABLE IF NOT EXISTS resources (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,              -- 'file', 'link', 'note'
    name TEXT NOT NULL,
    hash TEXT,                       -- only for file resources
    url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uniq_resources_hash
ON resources(hash)
WHERE hash IS NOT NULL;



CREATE TABLE IF NOT EXISTS resource_owners (
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    -- PRIMARY KEY (resource_id, user_id)
);

CREATE TABLE IF NOT EXISTS week_resources (
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    week_id UUID NOT NULL REFERENCES weeks(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
    -- PRIMARY KEY (resource_id, week_id)
);
