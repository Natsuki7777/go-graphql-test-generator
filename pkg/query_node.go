package pkg

import (
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
	Context      string
	Variables    string
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
	Name   string
	Type   string
	IsList bool
}

func QueryNode(tokens []string) (map[string]Query, map[string]Node) {
	qs := []Query{}
	nodes := []Node{}
	i := 0
	for i < len(tokens) {
		if tokens[i] == TYPE && tokens[i+1] == MUTATION { // type[i] Mutation { createUser(id:UUID ,input: CreateUserInput!): User! }
			i += 3 // type Mutation { createUser[i] (id:UUID ,input: CreateUserInput!): User! }
			for tokens[i] != BRACE_CLOSE {
				q := Query{}
				q.Type = MUTATION
				q.Name = tokens[i]
				i += 2 // type Mutation { createUser(id[i] :UUID ,input: CreateUserInput!): User! }
				for tokens[i] != PAREN_CLOSE {
					f := Field{}
					f.Name = tokens[i]
					i += 2 // type Mutation { createUser(id:UUID[i] ,input: CreateUserInput!): User! }
					f.Type = tokens[i]
					q.Inputs = append(q.Inputs, f)
					i++ // type Mutation { createUser(id:UUID ,[i] input : CreateUserInput!): User! }
					if tokens[i] == COMMA {
						i++ // type Mutation { createUser(id:UUID ,input[i] : CreateUserInput!): User! }
					}
				}
				i += 2 // type Mutation { createUser(id:UUID ,input: CreateUserInput! ): User![i] }
				if tokens[i] == SQUARE_OPEN {
					q.IsOutputList = true
					i++ // type Mutation { createUser(id:UUID ,input: CreateUserInput! ): [User][i] }
				}
				q.OutputTypes = append(q.OutputTypes, tokens[i])
				qs = append(qs, q)
				if q.IsOutputList {
					i++
				}
				i++
			}
			i++
		} else if tokens[i] == TYPE && tokens[i+1] == QUERY { // type[i] Query { user(id:UUID):[user]! }
			i += 3 // type Query { users[i] (id:UUID):[user]! }
			for tokens[i] != BRACE_CLOSE {
				q := Query{}
				q.Type = QUERY
				q.Name = tokens[i]
				i++ // type Query { users([i]id :UUID):[user]! }
				if tokens[i] == PAREN_OPEN {
					i++ // type Query { users(id[i] :UUID):[user]! }
					for tokens[i] != PAREN_CLOSE {
						f := Field{}
						f.Name = tokens[i]
						i += 2 // type Query { users(id:UUID[i] ):[user]! }
						f.Type = tokens[i]
						q.Inputs = append(q.Inputs, f)
						i++ // type Query { users(id:UUID[i] ):[user]! }
						if tokens[i] == COMMA {
							i++ // type Query { users(id:UUID, input[i] : CreateUserInput!): User! }
						}
					}
					i++ // type Query { users(id:UUID ):[i][user]! }
				}
				i++ // type Query { users(id:UUID ):[[i]user] }
				if tokens[i] == SQUARE_OPEN {
					q.IsOutputList = true
					i++ // type Query { users(id:UUID ):[[user[i]] }
				}
				q.OutputTypes = append(q.OutputTypes, tokens[i])
				qs = append(qs, q)
				if q.IsOutputList {
					i++ // type Query { users(id:UUID ):[user][i] }
				}
				i++ // type Query { users(id:UUID ):[user]! }[i]
			}
			i++
		} else if tokens[i] == INPUT { // input[i] UserInput { id: UUID! , name: String! }
			i++ // input UserInput[i] { id : UUID! , name: String! }
			n := Node{}
			n.Name = tokens[i]
			i += 2 // input UserInput { id[i] : UUID!  name: String! }
			for tokens[i] != BRACE_CLOSE {
				f := Field{}
				f.Name = tokens[i]
				i += 2 // input UserInput { id: UUID![i]  name: String! }
				f.Type = tokens[i]
				n.Fileds = append(n.Fileds, f)
				i++ // input UserInput { id: UUID!  name[i] : String! }
			}
			nodes = append(nodes, n)
			i++ // input UserInput { id: UUID!  name: String! }[i]
		} else { // type User { id: UUID! , scores: [Score!]! }
			i++ // type User[i] { id: UUID! , scores: [Score!]! }
			n := Node{}
			n.Name = tokens[i]
			i += 2 // type User { id[i] : UUID! , scores: [Score!]! }
			for tokens[i] != BRACE_CLOSE {
				f := Field{}
				f.Name = tokens[i]
				i += 2 // type User { id: UUID![i] , scores: [Score!]! }
				if tokens[i] == SQUARE_OPEN {
					f.IsList = true
					i++ // type User { id: UUID! , scores: [[i]Score!]! }
				}
				f.Type = tokens[i]
				n.Fileds = append(n.Fileds, f)
				i++ // type User { id: UUID! , scores: [Score!]![i] }
				if tokens[i] == SQUARE_CLOSE {
					i++ // type User { id: UUID! , scores: [Score!]! } [i]
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
			}
		}
		query_map[q.Name] = q
	}

	// Query struct
	//   Type: QUERY
	//   Name: createUser
	//   Inputs: [
	//     Field {
	//       Name: input
	//       Type: CreateUserInput
	//       IsList: false
	//   ]
	//   OutputTypes: [user]
	//	 Outputs: [
	//     Field {
	//       Name: id
	//       Type: UUID
	//       IsList: false
	//     },
	//   ]
	//   IsOutputList: true

	// Node struct
	//   Name: User
	//   Fileds: [
	//    {
	//     Name: id
	//     Type: UUID
	//		 IsList: false
	//    }
	//    {
	//     Name: scores
	//     Type: Score
	//		 IsList: true
	//    }
	//   ]
	return query_map, node_map
}

func GenerateContextVariables(q Query, node_map map[string]Node, count int, max_count int) (string, string) {
	context := ""
	variables := ""
	if q.Type == MUTATION {
		context += "mutation "
	} else if q.Type == QUERY {
		context += "query "
	}
	if q.Inputs != nil {
		context += "("
		variables += "{\n"
		query_input := FillQueryInput(q, node_map)
		variable_input := FillInputStruct(q.Inputs, node_map, count)
		variables += variable_input
		context += query_input
		context += ") {\n"
		variables += "}"
	} else {
		context += "{\n"
		context += "  " + q.Name + "{\n"
	}
	fileds := FillStruct(q.Outputs, node_map, count, max_count)
	context += fileds
	context += " }\n"

	return context, variables
}

func FillStruct(fields []Field, node_map map[string]Node, count int, max_depth int) string {
	temp := ""
	for _, f := range fields {
		if f.Type == TYPE_INT || f.Type == TYPE_FLOAT || f.Type == TYPE_STRING || f.Type == TYPE_BOOLEAN || f.Type == TYPE_UUID || f.Type == TYPE_DATETIME {
			temp += strings.Repeat("  ", count+2) + f.Name + "\n"
		} else {
			if count < max_depth {
				node := node_map[f.Type]
				temp += strings.Repeat("  ", count+2) + f.Name + " {\n"
				temp += FillStruct(node.Fileds, node_map, count+1, max_depth)
			}
		}
	}
	temp += strings.Repeat("  ", count+2) + "}\n"
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
			temp += strings.Repeat(" ", count*2) + ` "` + f.Name + `": "test"` + ",\n"
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

func FillQueryInput(q Query, node_map map[string]Node) string {
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

func FillQueryVariables(inputes []Field, node_map map[string]Node) string {
	temp := ""
	for _, f := range inputes {
		if f.Type == TYPE_INT {
			temp += "req.Var(\"" + f.Name + "\", 0)\n"
		} else if f.Type == TYPE_FLOAT {
			temp += "req.Var(\"" + f.Name + "\", 0.0)\n"
		} else if f.Type == TYPE_STRING {
			temp += "req.Var(\"" + f.Name + "\", \"test\")\n"
		} else if f.Type == TYPE_BOOLEAN {
			temp += "req.Var(\"" + f.Name + "\", false)\n"
		} else if f.Type == TYPE_UUID {
			temp += "req.Var(\"" + f.Name + "\", \"00000000-0000-0000-0000-000000000000\")\n"
		} else if f.Type == TYPE_DATETIME {
			temp += "req.Var(\"" + f.Name + "\", \"2000-01-01-00:00:00\")\n"
		} else {
			temp += "req.Var(\"input\", map[string]interface{}{\n"
			node := node_map[f.Type]
			temp += FillInputStruct(node.Fileds, node_map, 1)
			temp = temp[:len(temp)-1] + ",\n"
			temp += "})\n"
		}
	}
	return temp
}
