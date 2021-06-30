# redistore

Implementation of mini e-commerce app using Postgres, Redis, and Redisearch.
This repository is based on Hexagonal architecture.

## How to run project

Go step by step
- start postdres, redis and redisearch service instances by running:
        
        docker-compose up -d

- start application by:

        cd src/cmd && go run main.go initiator.go
