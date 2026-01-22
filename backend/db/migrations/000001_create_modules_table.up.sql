CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY,
    code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    department_name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS module_runs (
    id UUID PRIMARY KEY,
    module_id UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    year INT NOT NULL,
    semester TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (module_id, year, semester)
);

CREATE TABLE IF NOT EXISTS weeks (
    id UUID PRIMARY KEY, 
    module_run_id UUID NOT NULL REFERENCES module_runs(id) ON DELETE CASCADE,
    number INT NOT NULL
);

CREATE TABLE IF NOT EXISTS academic_terms (
    id UUID PRIMARY KEY, 
    year INT NOT NULL, 
    semester TEXT NOT NULL,
    is_active BOOLEAN NOT NULL
);

CREATE INDEX idx_module_runs_module_id
ON module_runs(module_id);
