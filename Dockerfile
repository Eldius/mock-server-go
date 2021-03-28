FROM golang:1.16-alpine3.13 as builder

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=1
ENV GOOS=linux

RUN apk add --no-cache git make build-base sqlite
RUN go build -v -a -ldflags '-extldflags "-static"' .

FROM golang:1.16-alpine3.13

WORKDIR /app
#RUN groupadd mocky && useradd -m -g mocky -l mocky
#RUN adduser -h /app -S -s /bin/ash -m -g mocky

COPY --chown=0:0 --from=builder /app/mock-server-go /app
COPY static /app/static
COPY mapper/samples/example_mapping_file.yml /app/mapping.yml
#USER mocky

ENTRYPOINT [ "/app/mock-server-go", "start", "/app/mapping.yml" ]
