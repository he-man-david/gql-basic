package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Tutorial struct {
	Id int
	Title string
	Author Author
	Comments []Comment
}

type Author struct {
	Name string
	Tutorials []int
}

type Comment struct {
	Body string
}

func populate() []Tutorial {
	author := &Author{Name: "David", Tutorials: []int{1}}
	tutorial := Tutorial{
		Id: 1,
		Title: "This is a tutorial",
		Author: *author,
		Comments: []Comment{
			{Body: "hello darkness my old friend..."},
		},
	}

	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial)

	return tutorials
}

func main() {
	fmt.Println("Graphql basic....")
	tutorials := populate()

	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"tutorials": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)

	var tutorialType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Tutorial",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"author": &graphql.Field{
					Type: authorType,
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(commentType),
				},
			},
		},
	)

	queryFields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type: tutorialType,
			Description: "get tutorial by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, tutorial := range tutorials {
						if int(tutorial.Id) == id {
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},
		"list": &graphql.Field{
			Type: graphql.NewList(tutorialType),
			Description: "Get full tutorial list",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: queryFields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new GraphQL Schema, ERR %v", err)
	}

	query := `
		{
			tutorial(id:1) {
				title
				author {
					name
					tutorials
				}
			}
		}
	`

	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to query Schema, ERR %v", r.Errors)
	}
	rJson, _ := json.Marshal(r)
	fmt.Printf("%s", &rJson)
}