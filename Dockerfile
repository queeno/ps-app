FROM golang@sha256:c0f8b5a435620fcf16601ad79b92a019e411bad0b94ad445bc84de0d09d8e0b8 AS builder
RUN apk update && apk add --no-cache git openssh ca-certificates tzdata musl-dev gcc build-base && update-ca-certificates
RUN adduser -D -g '' appuser
COPY *.go /ps-app/
COPY *.gohtml /ps-app/
COPY go.mod /ps-app/
WORKDIR /ps-app
ENV GOOS linux
ENV GOARCH amd64
RUN go test -v ./...
RUN go build -ldflags="-linkmode external -extldflags -static" -o /ps-app/bin/ps-app ./

FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /ps-app/index.gohtml /ps-app/index.gohtml
COPY --from=builder /ps-app/bin/ps-app /ps-app/bin/ps-app

USER appuser
EXPOSE 9292

WORKDIR /ps-app
ENTRYPOINT ["/ps-app/bin/ps-app"]
