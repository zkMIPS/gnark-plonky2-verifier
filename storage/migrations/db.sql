CREATE DATABASE IF NOT EXISTS zkm;

use zkm;

CREATE TABLE IF NOT EXISTS prover_job_queue
(
    id                  serial primary key,
    job_status          int       not null,
    job_priority        int       not null,
    job_type            text      not null,

    created_at          timestamp not null default now(),
    updated_by          text      not null,
    updated_at          timestamp not null default now(),

    proof_id            text      not null,
    computed_request_id text      not null,
    job_data            json      not null,
    INDEX proof_id_req_id(proof_id(255), computed_request_id(255)),
    INDEX job_status_updated_at(job_status, updated_at),
    INDEX job_type(job_type(255))
);

CREATE TABLE IF NOT EXISTS proofs
(
    id                  serial primary key,
    proof_id            text      not null,
    computed_request_id text      not null,
    proof               blob      not null,
    created_at          timestamp not null default now(),
    INDEX proof_id_req_id(proof_id(255), computed_request_id(255))
);

use mysql;

ALTER user 'root'@'localhost' IDENTIFIED BY '123456';
UPDATE user SET plugin='mysql_native_password' WHERE User='root';

FLUSH PRIVILEGES;