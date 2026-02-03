CREATE TABLE IF NOT EXISTS storage_objects(
    id  UUID PRIMARY KEY,
    hash TEXT NOT NULL, 
    url TEXT NOT NULL, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS resources(
    id UUID PRIMARY KEY, 
    type TEXT NOT NULL, 
    name TEXT NOT NULL, 
    storage_object_id UUID REFERENCES storage_objects(id), 
    external_url TEXT, 
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS resource_owners (
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (resource_id, user_id)
);

CREATE TABLE IF NOT EXISTS week_resources (
    resource_id UUID NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    week_id UUID NOT NULL REFERENCES weeks(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (resource_id, week_id)
);
select storage_objects that doestn have resources(id)

SELECT so.id FROM storage_objects so LEFT JOIN resources r ON so.id=r.storage_object_id WHERE r.id IS NULL;

