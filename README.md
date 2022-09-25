# Generate graphql, Go test files

This is a tool to generate graphql, Go test files.\
Mainly foucs on using this along with [gqlgen](https://github.com/99designs/gqlgen).

## Usage
Put your graphql schemas, queries, mutations, in `input_graphqls` folder. Then run `go run main.go` to generate graphql, Go test files. Output files will be in `outputs` folder.

## Dependencies
- [machinebox/graphql](https://github.com/machinebox/graphql)
  - Go client package for GraphQL APIs

## Todos
- [ ] Subscriptions
- [ ] Use more Template than parsing text
- [ ] Make it more configurable
- [ ] Command line arguments
- [ ] Use `go generate` instead of `go run main.go`


### My environment
go version go1.19.1 linux/amd64