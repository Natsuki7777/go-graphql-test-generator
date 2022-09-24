package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/Natsuki7777/go-graphql-test-generator/pkg"
)

func main() {
	targetDir := "input_graphqls"
	files, _ := ioutil.ReadDir(targetDir)
	targetString := ""
	for _, f := range files {
		fileName := f.Name()
		if strings.HasSuffix(fileName, "graphql") || strings.HasSuffix(fileName, "graphqls") {
			filePath := targetDir + "/" + fileName
			file, err := ioutil.ReadFile(filePath)
			if err != nil {
			}
			fileString := string(file)
			targetString += fileString
		}
	}

	tokens := pkg.TokenizeInput(targetString)
	query_map, node_map := pkg.QueryNode(tokens)

	result := ""

	for _, q := range query_map {
		temp := "\n"
		if q.Type == pkg.MUTATION {
			temp += "mutation ("
			input_string := pkg.FillMutationInput(q, node_map)
			temp += input_string + "){" + "\n"
			filed_string := pkg.FillStruct(q.Outputs, node_map, 1)
			temp += filed_string
			temp += "}\n"
			temp += "\n{\n"
			filled_string := pkg.FillInputStruct(q.Inputs, node_map, 0)
			temp += filled_string
			temp += "}\n"
		} else if q.Type == pkg.QUERY {
			if q.Inputs == nil {
				temp += "query {\n"
				temp += " " + q.Name + " {\n"
			} else {
				temp += "query ("
				for _, input := range q.Inputs {
					temp += "$" + input.Name + ": " + input.Type + "!"
					temp += ") {\n"
				}
				temp += " " + q.Name + "("
				for _, input := range q.Inputs {
					temp += input.Name + ": $" + input.Name
				}
				temp += ") {\n"
			}
			filed_string := pkg.FillStruct(q.Outputs, node_map, 1)
			temp += filed_string
			temp += " }\n"
			if q.Inputs != nil {
				temp += "\n{\n"
				filled_string := pkg.FillInputStruct(q.Inputs, node_map, 0)
				temp += filled_string
				temp += "}\n"
			}
		}
		result += temp
	}
	f, err := os.Create("output.graphql")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(result)
}
