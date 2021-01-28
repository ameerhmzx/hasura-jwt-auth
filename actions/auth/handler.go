package auth

import (
	"actions/utils"
	"crypto/rsa"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	jwtSignKey     *rsa.PrivateKey
	jwtVerifyKey   *rsa.PublicKey
	tokenExpTime   time.Duration
	refreshExpTime time.Duration
	adminSecret    string
	gqlServerUrl   string
	secretsPath    = "/usr/local/src/secrets"
	pemPath        = secretsPath + "/keys"
	wellKnownPath = secretsPath + "/.well-known"
)

func init() {
	// Load Env
	tokenExpTime = time.Duration(utils.LoadEnvInt("JWT_EXP_TIME", 5)) * time.Minute
	refreshExpTime = time.Duration(utils.LoadEnvInt("REF_EXP_DAYS", 7)*24) * time.Hour
	gqlServerUrl = utils.LoadEnvCritical("HASURA_GQL_API")
	adminSecret = utils.LoadEnvCritical("ADMIN_SECRET")

	verifyKeys()
}

func Handler(r *mux.Router) {
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/register", registerHandler)
	r.HandleFunc("/refresh", refreshHandler)
	r.PathPrefix("/.well-known/").Handler(http.StripPrefix("/.well-known/", http.FileServer(http.Dir(wellKnownPath))))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var actionPayload ActionPayloadRegister
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	result, err, errCode := registerNewUser(actionPayload.Input)

	if err != nil {
		errorObject := GraphQLError{
			Message: err.Error(),
		}
		errorBody, _ := json.Marshal(errorObject)
		w.WriteHeader(errCode)
		_, _ = w.Write(errorBody)
		return
	}

	data, _ := json.Marshal(result)
	_, _ = w.Write(data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var actionPayload ActionPayloadLogin
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	result, err, errCode := issueTokenFromCredentials(actionPayload.Input)

	if err != nil {
		errorObject := GraphQLError{
			Message: err.Error(),
		}
		errorBody, _ := json.Marshal(errorObject)
		w.WriteHeader(errCode)
		_, _ = w.Write(errorBody)
		return
	}

	data, _ := json.Marshal(result)
	_, _ = w.Write(data)
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	var actionPayload ActionPayloadRefresh
	err = json.Unmarshal(reqBody, &actionPayload)
	if err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	result, err, errCode := issueTokenFromRefreshToken(actionPayload.Input)

	if err != nil {
		errorObject := GraphQLError{
			Message: err.Error(),
		}
		errorBody, _ := json.Marshal(errorObject)
		w.WriteHeader(errCode)
		_, _ = w.Write(errorBody)
		return
	}

	data, _ := json.Marshal(result)
	_, _ = w.Write(data)
}
