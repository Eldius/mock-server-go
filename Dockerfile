FROM golang:1.18-alpine3.15 as builder

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=1

RUN apk add --no-cache git make build-base sqlite
RUN go build -v -a -ldflags '-extldflags "-static"' .
RUN chmod +x /app/mock-server-go

FROM alpine:3.15

EXPOSE 8080
EXPOSE 8081

WORKDIR /app

COPY --chown=0:0 --from=builder /app/mock-server-go /app
#COPY static /app/static
COPY mapper/samples/example_mapping_file.yml /app/mapping.yml

ENTRYPOINT [ "./mock-server-go", "start", "mapping.yml" ]
