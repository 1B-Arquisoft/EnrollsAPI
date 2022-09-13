package main

import (
	"SIA/InscripcionAPI/cmd/server"
	u "SIA/InscripcionAPI/cmd/utils"
)

func main() {
	server := server.Setup()

	server.Start(u.Get("API_URL"))
}
