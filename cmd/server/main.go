package main

import (
	"github.com/KnoblauchPilze/user-service/cmd/server/routes"
)

func main() {
	s := routes.NewServer()
	s.Start(":60000")
}
