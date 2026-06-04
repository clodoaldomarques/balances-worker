## build
FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY cmd/worker/main.go .
RUN go build -o /balances-worker

## deploy
FROM ubuntu AS base

WORKDIR /

COPY --from=build /balances-worker /balances-worker

#mysql access
ENV MYSQL_HOST=192.168.49.2
ENV MYSQL_PORT=30306
ENV MYSQL_DATABASE=balances
ENV MYSQL_USER=admin
ENV MYSQL_PASS=b4l4nc3s

#AWS 
ENV AWS_ADDRESS="http://192.168.49.2:4566"
ENV AWS_REGION="us-east-1"
ENV AWS_PROFILE="localstack"
ENV AWS_ID="test"
ENV AWS_PASS="test"

#SQS
ENV BALANCES_SQS_QUEUE="http://192.168.49.2:30566/000000000000/balances-sqs-queue"

EXPOSE 5001

ENTRYPOINT [ "/balances-worker" ]