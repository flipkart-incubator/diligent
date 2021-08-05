CREATE DATABASE diligent;
USE diligent;

CREATE TABLE baseline
(
    pk        varchar(32)    NOT NULL,
    uniq      varchar(32)    NOT NULL,
    small_grp varchar(32)    NOT NULL,
    large_grp varchar(32)    NOT NULL,
    fixed_val varchar(32)    NOT NULL,
    seq_num   int            NOT NULL,
    ts        timestamp      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    payload   varchar(10240) NOT NULL,
    PRIMARY KEY (pk)
);

CREATE TABLE experiment
(
    pk        varchar(32)    NOT NULL,
    uniq      varchar(32)    NOT NULL,
    small_grp varchar(32)    NOT NULL,
    large_grp varchar(32)    NOT NULL,
    fixed_val varchar(32)    NOT NULL,
    seq_num   int            NOT NULL,
    ts        timestamp      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    payload   varchar(10240) NOT NULL,
    PRIMARY KEY (pk),
    UNIQUE KEY idx_uniq (uniq),
    KEY idx_small_grp (small_grp),
    KEY idx_large_grp (large_grp),
    KEY idx_same (fixed_val),
    KEY idx_seq_num (seq_num),
    KEY idx_ts (ts)
);

