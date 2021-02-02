
start:
	go run main.go start mapper/samples/example_mapping_file.yml

test:
	go test ./... -cover
