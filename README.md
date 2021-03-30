# mock-server-go #

## build status ##

![Go](https://github.com/Eldius/mock-server-go/workflows/Go/badge.svg)

## dev snippets ##

Start server

```bash
go run main.go start mapper/samples/example_mapping_file.yml
```

curl #1

```bash
curl -i localhost:8080/v1/contract
```

curl #2

```bash
curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "123", "name": "test"}'
```

```bash
# benchmarking app on a Raspberry Pi 4 K3s cluster
wrk -c 20 -d 10m -H 'HeaderKey: HeaderValue' -H 'Cache-Control: no-cache' --timeout 3s -t 10 http://192.168.100.195:18080/v1/contract

```
## links ##

- [rogchap/v8go](https://github.com/rogchap/v8go)
