FROM golang:1.26-alpine AS builder

ARG VERSION

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w -X main.Version=${VERSION}" ./cmd/subscription

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=builder /build/subscription /bin/subscription

USER nobody
CMD ["/bin/subscription"]
