FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN go build -v main.go


FROM alpine:3.20

WORKDIR /api

ARG file

COPY --from=build /api/main ./
RUN mkdir config
COPY --from=build /api/config/${file}.yaml ./config

ENV file ${file}

CMD ./main -file ${file}
