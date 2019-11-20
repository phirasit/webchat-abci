FROM golang:latest

ARG ABCI_APP_PORT=8080
ARG ABCI_API_PORT=26657

ENV ABCI_APP_PORT $ABCI_APP_PORT

# build webchat ABCI
COPY . /go/src/webchatABCI
RUN go get -u github.com/tendermint/tendermint/...

RUN go install /go/src/webchatABCI
RUN go install github.com/tendermint/tendermint/cmd/tendermint

EXPOSE $ABCI_APP_PORT $ABCI_API_PORT

CMD ["sh", "/go/src/webchatABCI/start_abci.sh"]

