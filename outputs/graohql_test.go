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

func TestQueryuser(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	query ($id: UUID!) {
  user(id: $id) {
    id
    name
    group {
      id
      name
      createdAt
      updatedAt
      }
    role {
      id
      name
      createdAt
      updatedAt
      }
    createdAt
    updatedAt
    }
 }
`)

	req.Var("id", "00000000-0000-0000-0000-000000000000")

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestQueryuserGroup(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	query ($id: UUID!) {
  userGroup(id: $id) {
    id
    name
    createdAt
    updatedAt
    }
 }
`)

	req.Var("id", "00000000-0000-0000-0000-000000000000")

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestQueryuserGroups(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	query {
  userGroups{
    id
    name
    createdAt
    updatedAt
    }
 }
`)

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestQueryuserRole(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	query ($id: UUID!) {
  userRole(id: $id) {
    id
    name
    createdAt
    updatedAt
    }
 }
`)

	req.Var("id", "00000000-0000-0000-0000-000000000000")

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestQueryuserRoles(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	query {
  userRoles{
    id
    name
    createdAt
    updatedAt
    }
 }
`)

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestQueryusers(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	query {
  users{
    id
    name
    group {
      id
      name
      createdAt
      updatedAt
      }
    role {
      id
      name
      createdAt
      updatedAt
      }
    createdAt
    updatedAt
    }
 }
`)

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

//----------------------------MUTATION----------------------------

func TestMutaioncreateUser(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	mutation ($input: UserInput!) {
  createUser(input: $input) {
    id
    name
    group {
      id
      name
      createdAt
      updatedAt
      }
    role {
      id
      name
      createdAt
      updatedAt
      }
    createdAt
    updatedAt
    }
 }
`)

	req.Var("input", map[string]interface{}{
		"name":  "",
		"group": "00000000-0000-0000-0000-000000000000",
		"role":  "00000000-0000-0000-0000-000000000000",
	})

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestMutaioncreateUserGroup(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	mutation ($input: UserGroupInput!) {
  createUserGroup(input: $input) {
    id
    name
    createdAt
    updatedAt
    }
 }
`)

	req.Var("input", map[string]interface{}{
		"name": "",
	})

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestMutaioncreateUserRole(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	mutation ($input: UserRoleInput!) {
  createUserRole(input: $input) {
    id
    name
    createdAt
    updatedAt
    }
 }
`)

	req.Var("input", map[string]interface{}{
		"name": "",
	})

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestMutaionupdateUser(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	mutation ($id: UUID!,$input: UserInput!) {
  updateUser(id: $id,input: $input) {
    id
    name
    group {
      id
      name
      createdAt
      updatedAt
      }
    role {
      id
      name
      createdAt
      updatedAt
      }
    createdAt
    updatedAt
    }
 }
`)

	req.Var("id", "00000000-0000-0000-0000-000000000000")
	req.Var("input", map[string]interface{}{
		"name":  "",
		"group": "00000000-0000-0000-0000-000000000000",
		"role":  "00000000-0000-0000-0000-000000000000",
	})

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestMutaionupdateUserGroup(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	mutation ($id: UUID!,$input: UserGroupInput!) {
  updateUserGroup(id: $id,input: $input) {
    id
    name
    createdAt
    updatedAt
    }
 }
`)

	req.Var("id", "00000000-0000-0000-0000-000000000000")
	req.Var("input", map[string]interface{}{
		"name": "",
	})

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}

func TestMutaionupdateUserRole(t *testing.T) {
	client := graphql.NewClient(HOST)
	req := graphql.NewRequest(`
	mutation ($id: UUID!,$input: UserRoleInput!) {
  updateUserRole(id: $id,input: $input) {
    id
    name
    createdAt
    updatedAt
    }
 }
`)

	req.Var("id", "00000000-0000-0000-0000-000000000000")
	req.Var("input", map[string]interface{}{
		"name": "",
	})

	var respData interface{}
	if err := client.Run(context.Background(), req, &respData); err != nil {
		t.Fatal(err)
	}
	t.Log(respData)
}
