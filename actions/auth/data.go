package auth

import (
	"actions/utils"
	"context"
	"errors"
	"github.com/machinebox/graphql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func getUserInfo(username string) (string, []string, error, int) {
	client := graphql.NewClient(gqlServerUrl)
	req := graphql.NewRequest(`
		query($username: String!) {
		  auth_user(where: {username: {_eq: $username}}) {
			username
			password
			user_roles {
			  role
			  is_default
			}
		  }
		}
	`)
	req.Var("username", username)
	req.Header.Set("x-hasura-admin-secret", adminSecret)
	ctx := context.Background()

	//goland:noinspection GoSnakeCaseUsage
	var res struct {
		Auth_user []struct {
			Username   string
			Password   string
			User_roles []struct {
				Role       string
				Is_default bool
			}
		}
	}

	if err := client.Run(ctx, req, &res); err != nil {
		log.Fatal(err)
	}

	if len(res.Auth_user) != 1 {
		// there must be only one user with that username
		return "", nil, errors.New("username not registered"), http.StatusUnauthorized
	}

	roles := res.Auth_user[0].User_roles
	var rolesArr []string
	for _, role := range roles {
		if role.Is_default {
			rolesArr = append(rolesArr, "")
			copy(rolesArr[1:], rolesArr)
			rolesArr[0] = role.Role
		} else {
			rolesArr = append(rolesArr, role.Role)
		}
	}
	return res.Auth_user[0].Password, rolesArr, nil, http.StatusOK
}

func registerNewUser(args LoginArgs) (response Tokens, err error, errCode int) {
	mutation :=
		`mutation RegisterUser($username: String!, $password: String!, $role: String!) {
		  insert_auth_user_one(object: {username: $username, password: $password, user_roles: {data: {role: $role, is_default: true}}}) {
			username
		  }
		}`
	client := graphql.NewClient(gqlServerUrl)
	req := graphql.NewRequest(mutation)
	req.Header.Set("x-hasura-admin-secret", adminSecret)

	defRole := utils.LoadEnv("AUTHENTICATED_ROLE", "user")
	encPasswd, err := bcrypt.GenerateFromPassword([]byte(args.Credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return Tokens{}, errors.New("internal error"), http.StatusInternalServerError
	}

	req.Var("username", args.Credentials.Username)
	req.Var("password", string(encPasswd))
	req.Var("role", defRole)

	//goland:noinspection GoSnakeCaseUsage
	var res struct {
		Insert_auth_user_one struct {
			Username string
		}
	}

	if err := client.Run(context.Background(), req, &res); err != nil {
		return Tokens{}, errors.New("username already exist"), http.StatusConflict
	} else if res.Insert_auth_user_one.Username == args.Credentials.Username {
		return issueTokenFromCredentials(args)
	} else {
		return Tokens{}, errors.New("internal error"), http.StatusInternalServerError
	}
}
