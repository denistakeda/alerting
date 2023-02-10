CREATE TABLE metrics (
    id VARCHAR(256),
    mtype VARCHAR(10),
    value NUMERIC,
    delta BIGINT,
    UNIQUE (id, mtype)
);

CREATE UNIQUE INDEX id_mtype_index
ON metrics (id, mtype)
