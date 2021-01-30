package server

import "github.com/Eldius/mock-server-go/mapper"

func Start(port int, adminPort int, source string) {
	r := mapper.ImportMappingYaml(source)
	go StartAdminServer(adminPort, &r)
	StartMockServer(port, &r)
}
