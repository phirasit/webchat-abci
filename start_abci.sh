#!/bin/bash

# set port number if not set
ABCI_APP_PORT=${ABCI_APP_PORT:-8080}

# init tendermint node
tendermint init

# start tendermint node
tendermint node --proxy_app=localhost:${ABCI_APP_PORT} &

# start ABCI server
webchatABCI
