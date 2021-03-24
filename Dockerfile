FROM golang as builder

RUN mkdir /build

COPY . /build/

WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -o bin .

FROM golang

ARG PORT=5000
ARG API_HOST=""
ARG API_PORT=5000
ARG API_PUBLIC_PATH=/api/public
ARG API_TLS_DISABLED=false
ARG IMPORT_PATH=/api/seed.json
ARG REDIS_HOST=""
ARG REDIS_PORT=6379
ARG REDIS_PASSWORD=""

RUN mkdir /api

WORKDIR /build

COPY --from=builder /build/bin /api/
COPY seed.json /api/
COPY public /api/public

WORKDIR /api

LABEL   Name="Dinamicka Api"

#Run service
ENTRYPOINT ["./bin"]