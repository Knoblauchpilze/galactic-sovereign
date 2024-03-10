package main

import (
	"fmt"

	"github.com/KnoblauchPilze/user-service/cmd/server/routes"
)

const ENDPOINT = "users"
const PORT = uint16(60000)
const VERSION = 1

func main() {
	endpoint := fmt.Sprintf("/v%d/%s", VERSION, ENDPOINT)
	s := routes.NewServer(endpoint, PORT)

	s.Register(routes.UserRoutes())

	s.Start()
}
