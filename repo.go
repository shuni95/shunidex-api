package main

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

type PokemonTypeRepo struct {
	session neo4j.Session
	cache   map[string]PokemonTypeNode
}

type PokemonTypeNode struct {
	Name        string            `json:"name"`
	Translation string            `json:"translation"`
	Relations   []AgainstRelation `json:"relations"`
}

func (ptn *PokemonTypeNode) GetRelation(rivalType string) (float64, bool) {
	for _, relation := range ptn.Relations {
		if relation.Rival == rivalType {
			return relation.Effectiveness, true
		}
	}

	return -1.0, false
}

func (ptn PokemonTypeNode) IsStrongerAgainst() []PokemonTypeNode {
	var results []PokemonTypeNode

	for _, relation := range ptn.Relations {
		if relation.Effectiveness == 2.0 {
			results = append(results, typeRepo.cache[relation.Rival])
		}
	}

	return results
}

type AgainstRelation struct {
	Rival         string  `json:"rival"`
	Effectiveness float64 `json:"effectiveness"`
}

var typeRepo PokemonTypeRepo

func InitPokemonTypeRepo(session neo4j.Session) PokemonTypeRepo {
	typeRepo.session = session
	typeRepo.cache = typeRepo.GetAllPokemonType()

	return typeRepo
}

func (ptr PokemonTypeRepo) GetAllPokemonType() map[string]PokemonTypeNode {
	pokemonTypes := map[string]PokemonTypeNode{}

	_, err := ptr.session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, _ := tx.Run(`MATCH (pt:PokemonType) RETURN pt`, nil)

		for result.Next() {
			node := result.Record().Values[0].(dbtype.Node)
			name := node.Props["name"].(string)
			pokemonTypes[name] = PokemonTypeNode{
				Name:        name,
				Translation: node.Props["translation"].(string),
				Relations:   []AgainstRelation{},
			}
		}

		result, _ = tx.Run(`MATCH (mine:PokemonType)-[r:AGAINST]->(rival:PokemonType) 
		RETURN mine.name AS myType, rival.name AS rivalType, r.effectiveness AS effectiveness`, nil)

		for result.Next() {
			myType, _ := result.Record().Get("myType")
			rivalType, _ := result.Record().Get("rivalType")
			effectiveness, _ := result.Record().Get("effectiveness")

			// Add the against relation
			myTypeStr := myType.(string)
			if myEntry, ok := pokemonTypes[myTypeStr]; ok {
				myEntry.Relations = append(myEntry.Relations, AgainstRelation{
					Rival:         rivalType.(string),
					Effectiveness: effectiveness.(float64),
				})

				pokemonTypes[myTypeStr] = myEntry
			}
		}

		return nil, result.Err()
	})

	if err != nil {
		panic(err)
	}

	return pokemonTypes
}
