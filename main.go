package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

// Dog is a struct that contains the dog info.
type Dog struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Breed string `json:"breed"`
	Age   uint   `json:"age"`
}

// Dogs is an array of dogs.
var Dogs []Dog

func init() {
	dog1 := Dog{ID: 1, Name: "Ted", Breed: "Husky", Age: 3}
	dog2 := Dog{ID: 2, Name: "Bob", Breed: "Rottweiler", Age: 2}
	dog3 := Dog{ID: 3, Name: "Trap", Breed: "Dalmata", Age: 4}
	Dogs = append(Dogs, dog1, dog2, dog3)
}

var dogType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Dog",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"breed": &graphql.Field{
			Type: graphql.String,
		},
		"age": &graphql.Field{
			Type: graphql.Int,
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{

		"dog": &graphql.Field{
			Type:        dogType,
			Description: "Get single dog",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				idQuery, isOK := params.Args["id"].(uint)

				if isOK {
					for _, dog := range Dogs {
						if dog.ID == idQuery {
							return dog, nil
						}
					}
				}

				return Dog{}, nil
			},
		},

		"list": &graphql.Field{
			Type:        graphql.NewList(dogType),
			Description: "List of dogs",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return Dogs, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: nil,
})

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("The server is running on port: 8000")
	fmt.Println("Get single dog: curl -g 'http://localhost:8000/graphql?query={dog(id:1){id,name,breed,age}}'")
	fmt.Println("Get all dogs: curl -g 'http://localhost:8000/graphql?query={list{id, name, breed, age}}'")
	fmt.Println("Access the web app via browser at 'http://localhost:8000/graphql'")

	http.ListenAndServe(":8000", nil)
}
