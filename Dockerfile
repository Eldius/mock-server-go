FROM golang:1.24.5-alpine as builder

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=0

RUN apk add --no-cache git make build-base
RUN go build -o mock-server-go -v -a -ldflags '-extldflags "-static"' ./cmd/cli/
RUN ls -lha
RUN chmod +x /app/mock-server-go

FROM gcr.io/distroless/static

EXPOSE 8080
EXPOSE 8081

WORKDIR /app

COPY --chown=0:0 --from=builder /app/mock-server-go /app
COPY internal/mapper/samples/example_mapping_file.yml /app/mapping.yml

ENTRYPOINT [ "./mock-server-go", "start", "mapping.yml" ]
