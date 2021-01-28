package auth

import "github.com/dgrijalva/jwt-go"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
}

type RefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

// Hasura Args

type GraphQLError struct {
	Message string `json:"message"`
}

type Mutation struct {
	Login        *Tokens
	RegisterUser *Tokens
	RefreshToken *AccessToken
}

type LoginArgs struct {
	Credentials Credentials
}

type RefreshArgs struct {
	RefreshToken RefreshToken `json:"refreshToken"`
}

// JWT Claims

type AccessTokenClaims struct {
	HsClaims HasuraClaims `json:"https://hasura.io/jwt/claims"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	jwt.StandardClaims
	AuthClaims HasuraClaims `json:"auth_claims"`
}

type HasuraClaims struct {
	DefaultRole  string   `json:"x-hasura-default-role"`
	AllowedRoles []string `json:"x-hasura-allowed-roles"`
	Username     string   `json:"x-hasura-username"`
}

// Hasura Actions

type ActionPayloadLogin struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            LoginArgs              `json:"input"`
}

type ActionPayloadRegister struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            LoginArgs              `json:"input"`
}

type ActionPayloadRefresh struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
	Input            RefreshArgs           `json:"input"`
}
