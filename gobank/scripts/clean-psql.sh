#! /bin/bash

docker container rm gobank-postgres --force
docker run --name gobank-postgres -e POSTGRES_PASSWORD=gobank -p 5432:5432 -d postgres