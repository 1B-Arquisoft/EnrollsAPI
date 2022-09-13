package server

import (
	u "SIA/InscripcionAPI/cmd/utils"
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func neo4jConnection() neo4j.Session {
	driver, err := neo4j.NewDriver(u.Get("NEO4J_HOST"), neo4j.BasicAuth(u.Get("NEO4J_USER"), u.Get("NEO4J_PASSWORD"), ""))

	if err != nil {
		log.Fatal("Error Connecting to Neo4J DB")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	return session
}
