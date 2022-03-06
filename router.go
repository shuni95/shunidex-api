package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type PokemonTypeHandler struct {
	typeRepo *PokemonTypeRepo
}

type EvaluateTeamRequest struct {
	PokemonTypes [6][2]string
}

type EvaluateTeamResponse struct {
	WeakAgainstList   []string
	StrongAgainstList []string
}

func (ph *PokemonTypeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Println("Fetching all... (actually, just using cache)")

	payload, _ := json.Marshal(typeRepo.cache)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (ph *PokemonTypeHandler) GetType(w http.ResponseWriter, r *http.Request) {
	log.Println("Lookin for a specific type")
	vars := mux.Vars(r)

	typeNode := typeRepo.cache[vars["type"]]
	payload, _ := json.Marshal(typeNode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (ph *PokemonTypeHandler) EvaluateTeam(w http.ResponseWriter, r *http.Request) {
	log.Println("Evaluating")
	var teamRequest EvaluateTeamRequest
	var weaknessList []string

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&teamRequest)

	if err != nil {
		panic(err)
	}

	for _, pokemonTypes := range teamRequest.PokemonTypes {
		typeWeaknessList := ph.GetWeakness(pokemonTypes[0])
		for _, weakness := range typeWeaknessList {
			weaknessList = append(weaknessList, weakness.Name)
		}
		log.Println(pokemonTypes)
	}

	response := EvaluateTeamResponse{WeakAgainstList: weaknessList}
	payload, err := json.Marshal(response)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (ph *PokemonTypeHandler) GetWeakness(pokemonType string) []PokemonTypeNode {
	var results []PokemonTypeNode

	return results
}
