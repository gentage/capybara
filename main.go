package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gentage/capybara/pubsub"
	"github.com/gentage/capybara/resolver"
	"github.com/graph-gophers/graphql-go"

	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Self struct {
		Host string
		Port int
	}

	Redis struct {
		Host string
		Port int
	}
}

// TODO: handle errors
func loadConfig(filename string) (config Config) {
	bytes, _ := ioutil.ReadFile(filename)
	_ = yaml.Unmarshal(bytes, &config)
	return config
}

// TODO: handle errors
func loadFile(filename string) string {
	bytes, _ := ioutil.ReadFile(filename)
	return string(bytes)
}

// TODO: handle errors
func playground(config Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/playground.html")
		_ = t.Execute(w, fmt.Sprintf("http://%s:%d/", config.Self.Host, config.Self.Port))
	}
}

// TODO: handle errors
func gql(config Config) http.HandlerFunc {
	schema, _ := graphql.ParseSchema(
		loadFile("schema.gql"),
		resolver.NewResolver(pubsub.NewRedisClient(config.Redis.Host, config.Redis.Port)),
	)
	return graphqlws.NewHandlerFunc(schema, &relay.Handler{Schema: schema})
}

func main() {
	config := loadConfig("config.yaml")
	http.HandleFunc("/", gql(config))
	http.HandleFunc("/p", playground(config))

	formattedPort := fmt.Sprintf(":%d", config.Self.Port)
	log.Fatal(http.ListenAndServe(formattedPort, nil))
}
