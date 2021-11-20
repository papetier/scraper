# syntax=docker/dockerfile:1

FROM golang:1.17-alpine as builder

ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

# Install git
RUN apk add --no-cache git

WORKDIR /crawler

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go mod verify
RUN go install github.com/magefile/mage

COPY . .

RUN mage build:prod

FROM scratch as runner

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /crawler/build/crawler /crawler

USER appuser:appuser

ENTRYPOINT ["/crawler"]
