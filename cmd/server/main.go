package main

import (
	"github.com/KnoblauchPilze/template-go/cmd/server/routes"
)

func main() {
	e := routes.NewServer()

	e.Logger.Fatal(e.Start(":1323"))
}
