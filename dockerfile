FROM --platform=$BUILDPLATFORM golang:1.21 AS build-src
WORKDIR /src
ARG APP_VERSION="v0.0.0+unknown"
ARG TARGETOS="linux"
ARG TARGETARCH=""
#
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download

FROM --platform=$BUILDPLATFORM build-src AS build-server
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/member_activate-cli ./main.go

RUN ls -lah

FROM alpine:3.18
# ARG ENV
RUN apk --no-cache add ca-certificates tzdata
# ENV TZ=Asia/Bangkok
WORKDIR /app

COPY metadata-example.yaml /app/metadata_default.yaml
COPY scheduler-example.yaml /app/scheduler_default.yaml
COPY scheduler-worker-example.yaml /app/scheduler-worker_default.yaml

# Binary of application server
COPY --from=build-server /bin/member_activate-cli /app/

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN chown -R appuser:appgroup /app/*

USER appuser

ENTRYPOINT [ "/app/member_activate-cli" ]
