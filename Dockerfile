FROM golang:1.16-alpine as build_base
RUN apk add --no-cache git

WORKDIR /tmp/balancer

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./out

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build_base /tmp/balancer/out /app/balancer

ENTRYPOINT ["/app/balancer"]