package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"kalem/pubsub"
	"kalem/resolver"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Gql struct {
		Host string
		Port int
	}

	Redis struct {
		Host string
		Port int
	}
}

func loadConfig(filename string) (config Config) {
	bytes, _ := ioutil.ReadFile(filename)
	_ = yaml.Unmarshal(bytes, &config)
	return config
}

func loadFile(filename string) string {
	bytes, _ := ioutil.ReadFile(filename)
	return string(bytes)
}

func playground(endpoint string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/playground.html")
		_ = t.Execute(w, endpoint)
	}
}

func main() {
	config := loadConfig("config.yaml")

	schema := graphql.MustParseSchema(
		loadFile("schema.gql"),
		resolver.MakeResolver(pubsub.MakeRedisClient(config.Redis.Host, config.Redis.Port)),
	)

	gql := graphqlws.NewHandlerFunc(schema, &relay.Handler{Schema: schema})
	http.HandleFunc("/gql", gql)

	endpoint := fmt.Sprintf("http://%s:%d/gql", config.Gql.Host, config.Gql.Port)
	http.HandleFunc("/", playground(endpoint))

	formattedPort := fmt.Sprintf(":%d", config.Gql.Port)
	log.Fatal(http.ListenAndServe(formattedPort, nil))
}
