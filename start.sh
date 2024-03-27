#!/bin/bash

mysqld &
/usr/local/bin/snark_server -prover_cycle_time=15000 -log_level=4 > /app/server/server.log 2>&1