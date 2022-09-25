package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/Natsuki7777/go-graphql-test-generator/pkg"
)

const (
	TARGET_DIR  = "input_graphqls"
	MAX_DEPTH   = 6
	GRAPHOUTPUT = "outputs/graphql_test.graphql"
	GOOUTPUT    = "outputs/graohql_test.go"
)

func main() {
	targetDir := TARGET_DIR
	max_depth := MAX_DEPTH
	graph_output := GRAPHOUTPUT
	go_output := GOOUTPUT

	files, _ := ioutil.ReadDir(targetDir)
	targetString := ""
	for _, f := range files {
		fileName := f.Name()
		if strings.HasSuffix(fileName, "graphql") || strings.HasSuffix(fileName, "graphqls") {
			filePath := targetDir + "/" + fileName
			file, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Println(err)
			}
			fileString := string(file)
			targetString += fileString
		}
	}

	tokens := pkg.TokenizeInput(targetString)
	query_map, node_map := pkg.QueryNode(tokens)

	for k, q := range query_map {
		context, variables := pkg.GenerateContextVariables(q, node_map, 0, max_depth)
		q.Context = context
		q.Variables = variables
		query_map[k] = q
	}

	type Data struct {
		Query    []pkg.Query
		Mutation []pkg.Query
	}

	graphql_template, err := template.ParseFiles("template/graphql-template.graphqls")
	if err != nil {
		fmt.Println(err)
	}
	graphql_data := Data{}
	for _, q := range query_map {
		gq := pkg.Query{}
		gq.Name = q.Name
		gq.Context = q.Context
		if q.Variables != "" {
			gq.Variables = "# Variables\n#" + strings.Replace(q.Variables, "\n", "\n#", -1)
		}
		if q.Type == pkg.MUTATION {
			graphql_data.Mutation = append(graphql_data.Mutation, gq)
		} else if q.Type == pkg.QUERY {
			graphql_data.Query = append(graphql_data.Query, gq)
		}
	}
	sort.Slice(graphql_data.Query, func(i, j int) bool {
		return graphql_data.Query[i].Name < graphql_data.Query[j].Name
	})
	sort.Slice(graphql_data.Mutation, func(i, j int) bool {
		return graphql_data.Mutation[i].Name < graphql_data.Mutation[j].Name
	})
	f, err := os.Create(graph_output)
	if err != nil {
		panic(err)
	}
	graphql_template.Execute(f, graphql_data)
	defer f.Close()

	go_template, err := template.ParseFiles("template/go-template.go")
	if err != nil {
		panic(err)
	}
	go_data := Data{}
	for _, q := range query_map {
		gq := pkg.Query{}
		gq.Name = q.Name
		gq.Context = q.Context
		if q.Variables != "" {
			gq.Variables = pkg.FillQueryVariables(q.Inputs, node_map)
		}
		if q.Type == pkg.MUTATION {
			go_data.Mutation = append(go_data.Mutation, gq)
		} else if q.Type == pkg.QUERY {
			go_data.Query = append(go_data.Query, gq)
		}
	}
	sort.Slice(go_data.Query, func(i, j int) bool {
		return go_data.Query[i].Name < go_data.Query[j].Name
	})
	sort.Slice(go_data.Mutation, func(i, j int) bool {
		return go_data.Mutation[i].Name < go_data.Mutation[j].Name
	})
	f, err = os.Create(go_output)
	if err != nil {
		panic(err)
	}
	go_template.Execute(f, go_data)
	defer f.Close()
}
