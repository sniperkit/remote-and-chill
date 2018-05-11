package main

import (
	"encoding/json"
	"github.com/AndreasBackx/remote-and-chill/model"
	"github.com/AndreasBackx/remote-and-chill/resolver"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

func apiHandler(writer http.ResponseWriter, request *http.Request) {
	secretString := request.Header.Get("Authorization")
	secret, err := uuid.FromString(secretString)
	ctx := request.Context()

	if err == nil {
		ctx = model.Login(secret, ctx, resolver.Me)
	}

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		logrus.Error(err)
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	schema := graphql.MustParseSchema(SchemaString(), &resolver.Resolver{})
	response := schema.Exec(ctx, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		logrus.Error(err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(responseJSON)
}

func main() {
	http.HandleFunc("/", apiHandler)

	log.Fatal(http.ListenAndServe(":3000", nil))
}
