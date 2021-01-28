package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	jwtx "github.com/lestrrat-go/jwx/jwt"
	"io/ioutil"
	"os"
	"time"
)

var (
	privateKeyPath string
	publicKeyPath  string
)

type JWKS struct {
	Keys []jwk.Key `json:"keys"`
}

func verifyKeys() {
	privateKeyPath = pemPath + "/private.pem"
	publicKeyPath = pemPath + "/public.pem"

	if _, err := os.Stat(pemPath); os.IsNotExist(err) {
		err := os.MkdirAll(pemPath, os.ModePerm)
		if err != nil {
			fmt.Println("can't create directories!")
			panic(err)
		}
	}

	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		createKeys()
	} else if _, err := os.Stat(publicKeyPath); os.IsNotExist(err) {
		createKeys()
	}
	updateValues()
}

func updateValues() {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	checkError(err)
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	checkError(err)
	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	checkError(err)
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	checkError(err)
	jwtVerifyKey = verifyKey
	jwtSignKey = signKey

	key, err := jwk.New(verifyKey)
	checkError(err)
	_ = key.Set(jwtx.IssuedAtKey, time.Now().Unix())
	var jwks JWKS
	jwks.Keys = append(jwks.Keys, key)
	jwksStr, _ := json.MarshalIndent(jwks, "", "  ")

	if _, err := os.Stat(wellKnownPath); os.IsNotExist(err) {
		err := os.MkdirAll(wellKnownPath, os.ModePerm)
		if err != nil {
			fmt.Println("can't create directories!")
			panic(err)
		}
	}

	f, err := os.Create(wellKnownPath + "/jwks.json")
	checkError(err)
	defer f.Close()
	_, err = f.Write(jwksStr)
	_ = f.Sync()
	checkError(err)
}

func createKeys() {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	checkError(err)
	publicKey := key.PublicKey
	savePEMKey(privateKeyPath, key)
	savePublicPEMKey(publicKeyPath, publicKey)
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()
	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pubkey)
	checkError(err)
	var pemkey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}
	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()
	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
