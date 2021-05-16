-- Code generated with script in sql/migrations/templates DO NOT EDIT
BEGIN;

CREATE TABLE websitestatus(
    id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    url TEXT NOT NULL,
    regexp_pattern TEXT,
    regexp_match BOOLEAN,
    status_code SMALLINT,
    timetofirstbyte_ms SMALLINT,
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Create monthly partitions
CREATE TABLE websitestatus_2021_05 PARTITION OF websitestatus
    FOR VALUES FROM ('2021-05-01') TO ('2021-05-31');
CREATE TABLE websitestatus_2021_06 PARTITION OF websitestatus
    FOR VALUES FROM ('2021-06-01') TO ('2021-06-30');
CREATE TABLE websitestatus_2021_07 PARTITION OF websitestatus
    FOR VALUES FROM ('2021-07-01') TO ('2021-07-31');
CREATE TABLE websitestatus_2021_08 PARTITION OF websitestatus
    FOR VALUES FROM ('2021-08-01') TO ('2021-08-31');
CREATE TABLE websitestatus_2021_09 PARTITION OF websitestatus
    FOR VALUES FROM ('2021-09-01') TO ('2021-09-30');
CREATE TABLE websitestatus_2021_10 PARTITION OF websitestatus
    FOR VALUES FROM ('2021-10-01') TO ('2021-10-31');
CREATE TABLE websitestatus_2021_11 PARTITION OF websitestatus
    FOR VALUES FROM ('2021-11-01') TO ('2021-11-30');
CREATE TABLE websitestatus_2021_12 PARTITION OF websitestatus
    FOR VALUES FROM ('2021-12-01') TO ('2021-12-31');
CREATE TABLE websitestatus_2022_01 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-01-01') TO ('2022-01-31');
CREATE TABLE websitestatus_2022_02 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-02-01') TO ('2022-02-28');
CREATE TABLE websitestatus_2022_03 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-03-01') TO ('2022-03-31');
CREATE TABLE websitestatus_2022_04 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-04-01') TO ('2022-04-30');
CREATE TABLE websitestatus_2022_05 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-05-01') TO ('2022-05-31');
CREATE TABLE websitestatus_2022_06 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-06-01') TO ('2022-06-30');
CREATE TABLE websitestatus_2022_07 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-07-01') TO ('2022-07-31');
CREATE TABLE websitestatus_2022_08 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-08-01') TO ('2022-08-31');
CREATE TABLE websitestatus_2022_09 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-09-01') TO ('2022-09-30');
CREATE TABLE websitestatus_2022_10 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-10-01') TO ('2022-10-31');
CREATE TABLE websitestatus_2022_11 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-11-01') TO ('2022-11-30');
CREATE TABLE websitestatus_2022_12 PARTITION OF websitestatus
    FOR VALUES FROM ('2022-12-01') TO ('2022-12-31');
CREATE TABLE websitestatus_2023_01 PARTITION OF websitestatus
    FOR VALUES FROM ('2023-01-01') TO ('2023-01-31');
CREATE TABLE websitestatus_2023_02 PARTITION OF websitestatus
    FOR VALUES FROM ('2023-02-01') TO ('2023-02-28');
CREATE TABLE websitestatus_2023_03 PARTITION OF websitestatus
    FOR VALUES FROM ('2023-03-01') TO ('2023-03-31');
CREATE TABLE websitestatus_2023_04 PARTITION OF websitestatus
    FOR VALUES FROM ('2023-04-01') TO ('2023-04-30');
CREATE TABLE websitestatus_2023_05 PARTITION OF websitestatus
    FOR VALUES FROM ('2023-05-01') TO ('2023-05-31');

-- Index
CREATE INDEX ON websitestatus(created_at);

COMMIT;
