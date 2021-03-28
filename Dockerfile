FROM golang:1.16-alpine3.13 as builder

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=1
ENV GOOS=linux

RUN apk add --no-cache git make build-base sqlite
RUN go build -v -a -ldflags '-extldflags "-static"' .
RUN chmod +x /app/mock-server-go
FROM alpine:3.13

EXPOSE 8080
EXPOSE 8081

WORKDIR /app
#RUN groupadd mocky && useradd -m -g mocky -l mocky
#RUN adduser -h /app -S -s /bin/ash -m -g mocky

COPY --chown=0:0 --from=builder /app/mock-server-go /app
COPY static /app/static
COPY docker/entrypoint.sh /app/entrypoint.sh
COPY mapper/samples/example_mapping_file.yml /app/mapping.yml
#USER mocky

ENTRYPOINT [ "./mock-server-go", "start", "mapping.yml" ]
