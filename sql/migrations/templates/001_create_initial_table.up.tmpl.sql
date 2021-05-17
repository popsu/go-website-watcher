-- Code generated with script in sql/migrations/templates DO NOT EDIT
BEGIN;

CREATE TABLE websitestatus(
    id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    url TEXT NOT NULL,
    error TEXT,
    regexp_pattern TEXT,
    regexp_match BOOLEAN,
    status_code SMALLINT,
    timetofirstbyte_ms SMALLINT,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
{{- range $val := . }}
CREATE TABLE websitestatus_{{ $val.Format "2006_01"}} PARTITION OF websitestatus
    FOR VALUES FROM ('{{ $val.Format "2006-01-02" }}') TO ('{{ ($val.AddDate 0 1 -1 ).Format "2006-01-02" }}');
{{- end}}

-- Index
CREATE INDEX ON websitestatus(created_at);

COMMIT;
