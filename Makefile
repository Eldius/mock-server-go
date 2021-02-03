
start:
	go run main.go start mapper/samples/example_mapping_file.yml

test:
	-rm *.db
	-rm **/*.db
	go test ./... -cover

addtestroute:
	curl -i -XPOST localhost:8081/route -d @mapper/samples/example_route_request.json -H 'Content-Type: application/json'

getroutemapping:
	curl localhost:8081/route -H 'Accept: application/yaml'
