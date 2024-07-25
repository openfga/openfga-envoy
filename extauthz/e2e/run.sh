#!/bin/bash

set -e

DOCKER_COMPOSE='docker compose -f e2e/docker-compose.yaml'
YQ='yq'

STORE_ID=''
STORE_FILE='e2e/store.fga.yaml'

FGA_API_URL='http://localhost:18080'
TARGET_URL='http://localhost:8080'

which yq || (echo "yq is not installed. Please install it using make e2e-tools." && exit 1)
which fga || (echo "fga is not installed. Please install it make e2e-tools." && exit 1)

TMPDIR=$(mktemp -d)
MODEL=$TMPDIR/model.fga
$YQ '.model' $STORE_FILE > $MODEL

TUPLES=$TMPDIR/tuples.yaml
$YQ '.tuples' $STORE_FILE > $TUPLES

setup_fga_server() {
    $DOCKER_COMPOSE down
    echo "Setting FGA server."
    mkdir -p e2e/logs
    $DOCKER_COMPOSE up -d --build --remove-orphans openfga

    STORE_ID=$(fga store create --model $MODEL --api-url $FGA_API_URL | jq -rc '.store.id')
    echo "Created store with ID $STORE_ID"

    # TODO(jcchavezs): adds support for environment variable config to avoid this step
    STORE_ID=$STORE_ID envsubst < e2e/config.yaml.tmpl > e2e/config.yaml		
    $DOCKER_COMPOSE up -d --build envoy

    until [ "`docker inspect -f {{.State.Running}} envoy`"=="true" ]; do
        sleep 0.1;
    done;

    sleep 2
}

setup_fga_tuples() {
    echo "Writing FGA tuples."
    fga tuple write --store-id=$STORE_ID --file $TUPLES --api-url $FGA_API_URL | jq -er '.successful[0].object?' > /dev/null
}

failure () {
    cp $MODEL e2e/logs/model.fga
    cp $TUPLES e2e/logs/tuples.yaml
    $DOCKER_COMPOSE logs ext-authz > e2e/logs/ext-authz.log
    $DOCKER_COMPOSE logs envoy > e2e/logs/envoy.log
    $DOCKER_COMPOSE logs openfga > e2e/logs/openfga.log
    $DOCKER_COMPOSE down
    rm $MODEL
    rm $TUPLES
}

success() {
    $DOCKER_COMPOSE down
    rm $MODEL
    rm $TUPLES
}

do_call_and_expect() {
    expected_status_code=$1
    echo "Calling $TARGET_URL and expecting status code $expected_status_code."
    status_code=$(curl -s -o /dev/null -w "%{http_code}" $TARGET_URL)
    if [ "$status_code" -ne "$expected_status_code" ]; then
        echo "Expected status code $expected_status_code, got $status_code"
        failure
        exit 1
    fi
}

test_store() {
    fga model test --tests $STORE_FILE
}

run() {
    test_store
    setup_fga_server
    # Before setting the relationships
    do_call_and_expect 403
    setup_fga_tuples
    # After setting the relationships
    do_call_and_expect 200
    success
}

run