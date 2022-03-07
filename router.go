package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type PokemonTypeHandler struct {
	typeRepo *PokemonTypeRepo
}

type EvaluateTeamRequest struct {
	PokemonTypes [6][2]string
}

type EvaluateTeamResponse struct {
	WeakAgainst   map[string]EvaluateTeamEntryResponse `json:"weakAgainst"`
	StrongAgainst map[string]EvaluateTeamEntryResponse `json:"strongAgainst"`
}

type EvaluateTeamEntryResponse struct {
	Ocurrences  int    `json:"ocurrences"`
	Translation string `json:"translation"`
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
	weaknessesMap := map[string]EvaluateTeamEntryResponse{}
	strengthsMap := map[string]EvaluateTeamEntryResponse{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&teamRequest)

	if err != nil {
		panic(err)
	}

	for _, pokemonTypes := range teamRequest.PokemonTypes {
		typeOne := strings.ToLower(pokemonTypes[0])
		typeTwo := strings.ToLower(pokemonTypes[1])

		typeWeaknessList := ph.GetWeakness(typeOne, typeTwo)
		for _, weakness := range typeWeaknessList {
			if _, ok := weaknessesMap[weakness.Name]; !ok {
				weaknessesMap[weakness.Name] = EvaluateTeamEntryResponse{
					Ocurrences:  0,
					Translation: weakness.Translation,
				}
			}

			entry := weaknessesMap[weakness.Name]
			entry.Ocurrences = entry.Ocurrences + 1
			weaknessesMap[weakness.Name] = entry
		}

		// this should check the movements types, but for now we're using pokemon types
		for _, strongAgainst := range ph.typeRepo.cache[typeOne].IsStrongerAgainst() {
			if _, ok := strengthsMap[strongAgainst.Name]; !ok {
				strengthsMap[strongAgainst.Name] = EvaluateTeamEntryResponse{
					Ocurrences:  0,
					Translation: strongAgainst.Translation,
				}
			}

			entry := strengthsMap[strongAgainst.Name]
			entry.Ocurrences = entry.Ocurrences + 1
			strengthsMap[strongAgainst.Name] = entry
		}
	}

	response := EvaluateTeamResponse{WeakAgainst: weaknessesMap, StrongAgainst: strengthsMap}
	payload, err := json.Marshal(response)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func (ph *PokemonTypeHandler) GetWeakness(pokemonTypeOne, pokemonTypeTwo string) []PokemonTypeNode {
	var results []PokemonTypeNode

	for _, pokemonTypeCache := range ph.typeRepo.cache {
		effectivenessTypeOne, foundOne := pokemonTypeCache.GetRelation(pokemonTypeOne)
		if pokemonTypeTwo != "" {
			effectivenessTypeTwo, foundTwo := pokemonTypeCache.GetRelation(pokemonTypeTwo)

			// Check effectiveness combined
			if foundOne && foundTwo && (effectivenessTypeOne*effectivenessTypeTwo) >= 2.0 {
				results = append(results, pokemonTypeCache)
			} else if foundOne && !foundTwo && effectivenessTypeOne >= 2.0 {
				results = append(results, pokemonTypeCache)
			} else if foundTwo && !foundOne && effectivenessTypeTwo >= 2.0 {
				results = append(results, pokemonTypeCache)
			}
		} else {
			if foundOne && effectivenessTypeOne >= 2.0 {
				results = append(results, pokemonTypeCache)
			}
		}
	}

	return results
}
