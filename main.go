package main

import (
	"log"
	"net/http"

	_ "github.com/howtographql/graphql-go/db"

	"github.com/howtographql/graphql-go/resolvers"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"io/ioutil"
	"strings"
	"context"
)

var schema *graphql.Schema

func init() {
	schemaFile, err := ioutil.ReadFile("schema.graphqls")
	if err != nil {
		panic(err)
	}

	schema = graphql.MustParseSchema(string(schemaFile), &resolvers.Resolver{})
}

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next := &relay.Handler{Schema: schema}
		authorization := r.Header.Get("Authorization")
		token := strings.Replace(authorization, "Bearer ", "", 1)
		ctx := context.WithValue(r.Context(), "AuthorizationToken", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))

	log.Println("server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					headers: {
						'Accept': 'application/json',
						'Content-Type': 'application/json',
						'Authorization': 'Bearer 632788aa-781c-46e3-ad8d-825186c9c90b'
					},
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
