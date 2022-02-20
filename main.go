package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

type PokedexHandler struct {
	session neo4j.Session
}

func (ph PokedexHandler) GetAllPokemonHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("xd")
	_, err := ph.session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`MATCH (p:POKEMON) RETURN p`, nil)

		if err != nil {
			log.Fatal(err)
			panic(err)
		}

		for result.Next() {
			node := result.Record().Values[0].(dbtype.Node)
			if node.Props != nil {
				log.Println("Pokemon " + node.Props["name"].(string))
			}
		}

		return nil, result.Err()
	})

	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte("Shunidex!\n"))
}

func InitializeNeo4J() (neo4j.Driver, neo4j.Session) {
	log.Println("Connecting with neo4j with User " + os.Getenv("NEO4J_USER") + " and password " + os.Getenv("NEO4J_PASS"))
	uri := "neo4j://" + os.Getenv("NEO4J_HOST") + ":7687"
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
	ph := &PokedexHandler{session: session}

	// Routes consist of a path and a handler function.
	r.HandleFunc("/api/pokedex/all", ph.GetAllPokemonHandler)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", r))
}
