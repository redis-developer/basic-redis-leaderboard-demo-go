FROM golang as builder

WORKDIR /build

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

COPY . /build/

RUN CGO_ENABLED=0 GOOS=linux go build -o bin .

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
USER 65532:65532

ENV PORT=$PORT
ENV API_HOST=""
ENV API_PORT=8080
ENV API_PUBLIC_PATH=/api/public
ENV API_TLS_DISABLED=true
ENV IMPORT_PATH=/api/seed.json
ENV REDIS_HOST=""
ENV REDIS_PORT=6379
ENV REDIS_PASSWORD=""

WORKDIR /build

COPY --from=builder /build/bin /api/
COPY seed.json /api/
COPY public /api/public

WORKDIR /api

LABEL   Name="Leaderboard Api"

#Run service
ENTRYPOINT ["./bin"]
