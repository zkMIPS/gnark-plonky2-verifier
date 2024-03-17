CREATE DATABASE IF NOT EXISTS zkm;

use zkm;

DROP TABLE IF EXISTS prover_job_queue;
CREATE TABLE prover_job_queue
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
    INDEX proof_id_req_id(proof_id(255), computed_request_id(255))
);

DROP TABLE IF EXISTS proofs;
CREATE TABLE proofs
(
    proof_id            text      not null,
    computed_request_id text      not null,
    proof               blob      not null,
    created_at          timestamp not null default now(),
    INDEX proof_id_req_id(proof_id(255), computed_request_id(255))
);