package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func genAccessToken(username string, roles []string, defaultRole string) (string, error) {
	claims := &AccessTokenClaims{
		HsClaims: HasuraClaims{
			Username:     username,
			AllowedRoles: roles,
			DefaultRole:  defaultRole,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if tokenString, err := token.SignedString(jwtSignKey); err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func genRefreshToken(username string, roles []string) (string, error) {
	claims := &RefreshTokenClaims{
		AuthClaims: HasuraClaims {
			Username:     username,
			AllowedRoles: roles,
			DefaultRole:  roles[0],
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshExpTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	if tokenString, err := token.SignedString(jwtSignKey); err != nil {
		return "", err
	} else {
		return tokenString, nil
	}
}

func issueTokenFromCredentials(args LoginArgs) (response Tokens, err error, httpCode int) {
	response = Tokens{}

	errBadRequest := errors.New("bad request")
	errBadPassword := errors.New("wrong password")
	errInternal := errors.New("internal error")

	credentials := args.Credentials
	if credentials.Username != "" && credentials.Password != "" {
		encPasswd, roles, err, statusCode := getUserInfo(credentials.Username)
		if err != nil {
			return response, err, statusCode
		}
		if err := bcrypt.CompareHashAndPassword([]byte(encPasswd), []byte(credentials.Password)); err == nil {
			// Generate token
			response.AccessToken, err = genAccessToken(credentials.Username, roles, roles[0])
			if err != nil {
				return response, errInternal, http.StatusInternalServerError
			}
			response.RefreshToken, err = genRefreshToken(credentials.Username, roles)
			if err != nil {
				return response, errInternal, http.StatusInternalServerError
			}
			// generate refresh Tokens
			return response, nil, http.StatusOK
		} else {
			return response, errBadPassword, http.StatusUnauthorized
		}
	} else {
		return response, errBadRequest, http.StatusBadRequest
	}
}

func issueTokenFromRefreshToken(args RefreshArgs) (response AccessToken, err error, httpCode int) {
	response = AccessToken{}
	claims := RefreshTokenClaims{}

	errInternal := errors.New("internal error")
	errRefToken := errors.New("refresh token expired")

	token, err := jwt.ParseWithClaims(args.RefreshToken.RefreshToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return jwtVerifyKey, nil
	})

	if err != nil {
		return response, errInternal, http.StatusInternalServerError
	}

	if token.Valid {
		response.AccessToken, err = genAccessToken(claims.AuthClaims.Username, claims.AuthClaims.AllowedRoles, claims.AuthClaims.DefaultRole)
		if err != nil {
			return response, errInternal, http.StatusInternalServerError
		}
		return response, nil, http.StatusOK
	} else {
		return response, errRefToken, http.StatusUnauthorized
	}

}
