# Golang Boilerplate for Integration Test

We have containers. It provides docker images to run integration tests.

Test suites provides the suite to test our repository implementation.

### Run `docker` then `go test ./...` to run integration tests

## Neo4j

`neo4j:4.3-community` image used for integration test

## Couchbase

`docker.io/trendyoltech/couchbase-testcontainer:6.5.1` image used for integration test

This image provides environment variables to set username, password and bucket etc.
So, we can integrate easily to our project.
