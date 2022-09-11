package main

import (
	"SIA/InscripcionAPI/cmd/server"
)

func main() {
	server := server.Setup()

	server.Start("0.0.0.0:8080")
}
