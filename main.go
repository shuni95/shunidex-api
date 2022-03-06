package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func InitializeNeo4J() (neo4j.Driver, neo4j.Session) {
	log.Println("Connecting with neo4j with User " + os.Getenv("NEO4J_USER") + " and password " + os.Getenv("NEO4J_PASS"))
	uri := os.Getenv("NEO4J_URI")
	auth := neo4j.BasicAuth(os.Getenv("NEO4J_USER"), os.Getenv("NEO4J_PASS"), "")
	driver, err := neo4j.NewDriver(uri, auth)

	if err != nil {
		log.Fatal(err)
	}

	session := driver.NewSession(neo4j.SessionConfig{})

	return driver, session
}

func main() {
	log.Println("Starting Shunidex API...")
	// Wait 5 seconds for Neo4J
	time.Sleep(5 * time.Second)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	driver, session := InitializeNeo4J()
	defer driver.Close()
	defer session.Close()

	r := mux.NewRouter()
	ptr := InitPokemonTypeRepo(session)
	pth := &PokemonTypeHandler{typeRepo: &ptr}

	// Routes consist of a path and a handler function.
	r.HandleFunc("/api/types", pth.GetAll)
	r.HandleFunc("/api/types/{type}", pth.GetType)
	r.HandleFunc("/api/types/evaluateTeam", pth.EvaluateTeam)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
