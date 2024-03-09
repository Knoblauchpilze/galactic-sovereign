package main

import (
	"github.com/KnoblauchPilze/user-service/cmd/server/routes"
)

func main() {
	s := routes.NewServer("/users", uint16(60000))

	s.Register(routes.UserRoutes())

	s.Start()
}
