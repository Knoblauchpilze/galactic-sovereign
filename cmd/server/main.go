package main

import (
	"fmt"

	"github.com/KnoblauchPilze/user-service/cmd/server/routes"
	"github.com/KnoblauchPilze/user-service/cmd/server/server"
)

const ENDPOINT = "users"
const PORT = uint16(60000)
const VERSION = 1

func main() {
	endpoint := fmt.Sprintf("/v%d/%s", VERSION, ENDPOINT)
	s := server.New(endpoint, PORT)

	s.Register(routes.UserRoutes())

	s.Start()
}
