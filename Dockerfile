FROM golang as builder

RUN mkdir /build

COPY . /build/

WORKDIR /build

RUN CGO_ENABLED=0 GOOS=linux go build -o bin .

FROM golang

ENV PORT=$PORT
ENV API_HOST=""
ENV API_PORT=5000
ENV API_PUBLIC_PATH=/api/public
ENV API_TLS_DISABLED=false
ENV IMPORT_PATH=/api/seed.json
ENV REDIS_HOST=""
ENV REDIS_PORT=6379
ENV REDIS_PASSWORD=""

RUN mkdir /api

WORKDIR /build

COPY --from=builder /build/bin /api/
COPY seed.json /api/
COPY public /api/public

WORKDIR /api

LABEL   Name="Dinamicka Api"

#Run service
ENTRYPOINT ["./bin"]