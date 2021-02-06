
clean:
	-rm *.db*
	-rm **/*.db*

start:
	go run main.go start mapper/samples/example_mapping_file.yml

test: clean
	go test ./... -cover

addtestroute:
	curl -i -XPOST localhost:8081/route -d @mapper/samples/example_route_request.json -H 'Content-Type: application/json'

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
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "126", "name": "test3"}'
	curl -i localhost:8080/v1/contract
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "127", "name": "test4"}'
	curl -i localhost:8080/v1/contract
	curl -i -XPOST http://localhost:8080/v1/contract -d '{"id": "128", "name": "test5"}'
