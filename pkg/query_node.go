package pkg

import (
	"fmt"
	"strings"
)

const (
	MUTATION = "Mutation"
	QUERY    = "Query"
	NODE     = "Node"
	INPUT    = "input"
	TYPE     = "type"
)

const (
	TYPE_STRING   = "String"
	TYPE_INT      = "Int"
	TYPE_UUID     = "UUID"
	TYPE_DATETIME = "DateTime"
	TYPE_BOOLEAN  = "Boolean"
	TYPE_FLOAT    = "Float"
)

type Query struct {
	Type         string
	Name         string
	InputTypes   []string
	Inputs       []Field
	OutputTypes  []string
	Outputs      []Field
	IsOutputList bool
}

type Node struct {
	Name   string
	Fileds []Field
}

type Field struct {
	Name string
	Type string
}

func QueryNode(tokens []string) (map[string]Query, map[string]Node) {
	qs := []Query{}
	nodes := []Node{}
	i := 0
	for i < len(tokens) {
		if tokens[i] == TYPE && tokens[i+1] == MUTATION { // type Mutation
			i += 3
			for tokens[i] != BRACE_CLOSE {
				q := Query{}
				q.Type = MUTATION
				q.Name = tokens[i]
				i += 2
				for tokens[i] != PAREN_CLOSE {
					if tokens[i] == INPUT {
						i += 2
						q.InputTypes = append(q.InputTypes, tokens[i])
					} else {
						f := Field{}
						f.Name = tokens[i]
						i += 2
						f.Type = tokens[i]
						q.Inputs = append(q.Inputs, f)
					}
					i++
					if tokens[i] == COMMA {
						i++
					}
				}
				i += 2
				if tokens[i] == SQUARE_OPEN {
					fmt.Println("list")
					q.IsOutputList = true
					i++
				}
				q.OutputTypes = append(q.OutputTypes, tokens[i])
				qs = append(qs, q)
				if q.IsOutputList {
					i++
				}
				i++
			}
			i++
		} else if tokens[i] == TYPE && tokens[i+1] == QUERY { // type Query
			i += 3
			for tokens[i] != BRACE_CLOSE {
				q := Query{}
				q.Type = QUERY
				q.Name = tokens[i]
				i++
				if tokens[i] == PAREN_OPEN {
					i++
					for tokens[i] != PAREN_CLOSE {
						f := Field{}
						f.Name = tokens[i]
						i += 2
						f.Type = tokens[i]
						q.Inputs = append(q.Inputs, f)
						i++
						if tokens[i] == COMMA {
							i++
						}
					}
					i++
				}
				if tokens[i] == COLON {
					i++
				}
				if tokens[i] == SQUARE_OPEN {
					q.IsOutputList = true
					i++
				}
				q.OutputTypes = append(q.OutputTypes, tokens[i])
				qs = append(qs, q)
				if q.IsOutputList {
					i++
				}
				i++
			}
			i++
		} else if tokens[i] == INPUT { // input xxx
			i++
			n := Node{}
			n.Name = tokens[i]
			i += 2
			for tokens[i] != BRACE_CLOSE {
				f := Field{}
				f.Name = tokens[i]
				i += 2
				f.Type = tokens[i]
				n.Fileds = append(n.Fileds, f)
				i++
			}
			nodes = append(nodes, n)
			i++
		} else { // type xxx
			i++
			n := Node{}
			n.Name = tokens[i]
			i += 2
			for tokens[i] != BRACE_CLOSE {
				f := Field{}
				f.Name = tokens[i]
				i += 2
				if tokens[i] == SQUARE_OPEN {
					i++
				}
				f.Type = tokens[i]
				n.Fileds = append(n.Fileds, f)
				i++
				if tokens[i] == SQUARE_CLOSE {
					i++
				}
			}
			nodes = append(nodes, n)
			i++
		}
	}

	node_map := map[string]Node{}
	for _, n := range nodes {
		node_map[n.Name] = n
	}
	query_map := map[string]Query{}

	for _, q := range qs {
		for _, t := range q.OutputTypes {
			if _, ok := node_map[t]; ok {
				q.Outputs = append(q.Outputs, node_map[t].Fileds...)
				for _, i := range q.InputTypes {
					q.Inputs = append(q.Inputs, Field{Name: i, Type: i})
				}
			}
		}
		query_map[q.Name] = q
	}

	return query_map, node_map
}

func FillStruct(fields []Field, node_map map[string]Node, count int) string {
	temp := ""
	for _, f := range fields {
		if f.Type == TYPE_INT || f.Type == TYPE_FLOAT || f.Type == TYPE_STRING || f.Type == TYPE_BOOLEAN || f.Type == TYPE_UUID || f.Type == TYPE_DATETIME {
			temp += strings.Repeat(" ", count*2) + f.Name + "\n"
		} else {
			node := node_map[f.Type]
			temp += strings.Repeat(" ", count*2) + f.Name + " {\n"
			if count < 5 {
				temp += FillStruct(node.Fileds, node_map, count*2)
			}
		}
	}
	temp += strings.Repeat(" ", count*2) + "}\n"
	return temp
}

func FillInputStruct(fields []Field, node_map map[string]Node, count int) string {
	temp := ""
	for _, f := range fields {
		if f.Type == TYPE_INT {
			temp += strings.Repeat(" ", count*2) + `"` + f.Name + `": 0` + ",\n"
		} else if f.Type == TYPE_FLOAT {
			temp += strings.Repeat(" ", count*2) + ` "` + f.Name + `": 0.0` + ",\n"
		} else if f.Type == TYPE_STRING {
			temp += strings.Repeat(" ", count*2) + ` "` + f.Name + `": "NULL"` + ",\n"
		} else if f.Type == TYPE_BOOLEAN {
			temp += strings.Repeat(" ", count*2) + ` "` + f.Name + `": false` + ",\n"
		} else if f.Type == TYPE_UUID {
			temp += strings.Repeat(" ", count*2) + ` "` + f.Name + `": "00000000-0000-0000-0000-000000000000"` + ",\n"
		} else if f.Type == TYPE_DATETIME {
			temp += strings.Repeat(" ", count*2) + ` "` + f.Name + `": "2000-01-01-00:00:00"` + ",\n"
		} else {
			temp += strings.Repeat(" ", count*2) + ` "input": {` + "\n"
			node := node_map[f.Type]
			temp += FillInputStruct(node.Fileds, node_map, count+1)
			temp += strings.Repeat(" ", count*2) + ` }` + ",\n"
		}
	}
	if temp[len(temp)-2:] == ",\n" {
		temp = temp[:len(temp)-2] + "\n"
	}
	return temp
}

func FillMutationInput(q Query, node_map map[string]Node) string {
	temp := ""
	for _, f := range q.Inputs {
		if f.Type == TYPE_INT {
			temp += "$" + f.Name + ": Int!"
		} else if f.Type == TYPE_FLOAT {
			temp += "$" + f.Name + ": Float!"
		} else if f.Type == TYPE_STRING {
			temp += "$" + f.Name + ": String!"
		} else if f.Type == TYPE_BOOLEAN {
			temp += "$" + f.Name + ": Boolean!"
		} else if f.Type == TYPE_UUID {
			temp += "$" + f.Name + ": UUID!"
		} else if f.Type == TYPE_DATETIME {
			temp += "$" + f.Name + ": DateTime!"
		} else {
			temp += "$input" + ": " + f.Type + "!"
		}
		if f != q.Inputs[len(q.Inputs)-1] {
			temp += ","
		}
	}
	temp += ") {\n"
	temp += "  " + q.Name + "("
	for _, f := range q.Inputs {
		if f.Type == TYPE_INT {
			temp += f.Name + ": $" + f.Name
		} else if f.Type == TYPE_FLOAT {
			temp += f.Name + ": $" + f.Name
		} else if f.Type == TYPE_STRING {
			temp += f.Name + ": $" + f.Name
		} else if f.Type == TYPE_BOOLEAN {
			temp += f.Name + ": $" + f.Name
		} else if f.Type == TYPE_UUID {
			temp += f.Name + ": $" + f.Name
		} else if f.Type == TYPE_DATETIME {
			temp += f.Name + ": $" + f.Name
		} else {
			temp += "input" + ": $" + "input"
		}
		if f != q.Inputs[len(q.Inputs)-1] {
			temp += ","
		}
	}
	return temp
}
