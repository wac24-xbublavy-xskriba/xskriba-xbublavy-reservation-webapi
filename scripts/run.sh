#!/bin/bash

command=$1

if [ -z "$command" ]; then
    command="start"
fi

ProjectRoot="$(dirname "$(readlink -f "$0")")/.."

export RESERVATION_API_ENVIRONMENT="Development"
export RESERVATION_API_PORT="8080"
export RESERVATION_API_MONGODB_USERNAME="root"
export RESERVATION_API_MONGODB_PASSWORD="neUhaDnes"

mongo() {
    docker-compose --file "${ProjectRoot}/deployments/docker-compose/compose.yaml" "$@"
}

shutdown_mongo() {
    mongo down
}

trap 'shutdown_mongo' EXIT

case "$command" in
    "openapi")
        docker run --rm -ti -v "${ProjectRoot}:/local" openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
        ;;
    "start")
        mongo up --detach
        go run "${ProjectRoot}/cmd/reservation-api-service"
        ;;
    "mongo")
        mongo up
        ;;
    *)
        echo "Unknown command: $command"
        exit 1
        ;;
esac
