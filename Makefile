
clean:
	-rm *.db*
	-rm **/*.db*

start:
	go run ./cmd/cli/main.go start mapper/samples/example_mapping_file.yml

test: clean
	go test ./... -cover

addtestroute:
	curl -i -XPOST localhost:8081/route -d @internal/mapper/samples/example_route_request.json -H 'Content-Type: application/json'

getroutemapping:
	curl localhost:8081/route -H 'Accept: application/yaml'

benchmark: clean
	go test -run=Bench -bench=. ./...

makerequests:
	curl -i localhost:8080/v1/contract
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "123", "name": "test0"}'
	curl -i localhost:8080/v1/contract
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "124", "name": "test1"}'
	curl -i localhost:8080/v1/contract
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "125", "name": "test2"}'
	curl -i localhost:8080/v1/contract
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "126", "name": "test3"}' -H 'Content-Type: application/json'
	curl -i localhost:8080/v1/contract
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "127", "name": "test4"}' -H 'Content-Type: application/json'
	curl -i localhost:8080/v1/contract
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "128", "name": "test5"}' -H 'Content-Type: application/json'
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "128", "name": "test5"}' -H 'Content-Type: application/json'
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "128", "name": "test5"}' -H 'Content-Type: application/json'
	curl -i -XPOST http://localhost:8080/v2/test -d '{"id": "128", "contract": 123450, "name": "test0"}' -H 'Content-Type: application/json'
	curl -i -XPOST http://localhost:8080/v2/test -d '{"id": "128", "contract": 123451, "name": "test1"}' -H 'Content-Type: application/json'
	curl -i -XPOST http://localhost:8080/v2/test -d '{"id": "128", "contract": 123452, "name": "test2"}' -H 'Content-Type: application/json'
	curl -i -XPOST http://localhost:8080/v2/test -d '{"id": "128", "contract": 123453, "name": "test3"}' -H 'Content-Type: application/json'

dockerbuild:
	docker build \
		-t eldius/mock-server-go \
		.
	docker tag eldius/mock-server-go eldius/mock-server-go:$(shell git rev-parse --short HEAD)

dockerbuildarm:
	docker build \
		-t eldius/mock-server-go-armhf \
		-f Dockerfile.armhf \
		.
	docker tag eldius/mock-server-go-armhf eldius/mock-server-go-armhf:$(shell git rev-parse --short HEAD)

dockerpush: dockerbuild
	docker push eldius/mock-server-go:latest
	docker push eldius/mock-server-go:$(shell git rev-parse --short HEAD)

dockerpusharm: dockerbuildarm
	docker push eldius/mock-server-go-armhf:latest
	docker push eldius/mock-server-go-armhf:$(shell git rev-parse --short HEAD)

dockerrun: dockerbuild
	docker run -it --rm --name mocky -p 8080:8080 -p 8081:8081 eldius/mock-server-go:latest

dockermulti:
	docker buildx build \
		--push \
		--platform linux/arm/v7,linux/arm64/v8,linux/amd64 \
		--tag eldius/mock-server-go-multi:latest .

request/request.pb.go:
	protoc -I=$(shell pwd) --go_out=$(shell pwd) $(shell pwd)/request/proto/request.proto

lint:
	revive -formatter stylish ./...
