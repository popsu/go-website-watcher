INSERT INTO websitestatus (
    id,
    created_at,
    url,
    regexp_pattern,
    regexp_match,
    status_code,
    timetofirstbyte_ms
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
