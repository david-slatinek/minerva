FROM golang:1.23.1-alpine3.20 AS build

WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
RUN go build -v main.go


FROM alpine:3.20

RUN apk --no-cache add curl

WORKDIR /api

ARG FILE

COPY --from=build /api/main ./
RUN mkdir config
COPY --from=build /api/config/${FILE}.yaml ./config

ARG VERSION
ENV VERSION=${VERSION}

ENV file=${FILE}

CMD ["sh", "-c", "./main -file ${file}"]

HEALTHCHECK --interval=10s --timeout=5s --retries=3 CMD curl -f http://localhost:8080/health
