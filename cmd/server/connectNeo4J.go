package server

import (
	"log"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func neo4jConnection() neo4j.Session {
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", `test`, ""))

	if err != nil {
		log.Fatal("Error Connecting to Neo4J DB")
	}

	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "inscripciones"})
	return session
}
