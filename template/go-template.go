package test

import (
	"context"
	"testing"

	"github.com/machinebox/graphql"
)

const (
	HOST = "http://localhost:8080/query"
)

//----------------------------QUERY----------------------------
{{range .Query}}
func TestQuery{{.Name}}(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	{{.Context}}`)
	{{if ne .Variables  ""}}
	{{.Variables}}
	{{end}}
	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}
{{end}}


//----------------------------MUTATION----------------------------
{{range .Mutation}}
func TestMutaion{{.Name}}(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	{{.Context}}`)
	{{if ne .Variables  ""}}
	{{.Variables}}
	{{end}}
	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}
{{end}}
