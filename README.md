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

## links ##

- [rogchap/v8go](https://github.com/rogchap/v8go)
