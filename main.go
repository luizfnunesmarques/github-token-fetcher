package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v45/github"
)

type installationTokenResponse struct {
	Token        string                         `json:"token"`
	ExpiresAt    time.Time                      `json:"expires_at"`
	Permissions  github.InstallationPermissions `json:"permissions,omitempty"`
	Repositories []github.Repository            `json:"repositories,omitempty"`
}

func main() {
	installationID := flag.String("installationID", "", "The Application's installation ID.")
	applicationID := flag.String("applicationID", "", "Original Application ID.")
	privateKeyFilePath := flag.String("privateKeyFilePath", "", "Path where to find the pem file.")

	flag.Parse()

	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() == "" {
			log.Fatalf("Flag not set: %s", f.Name)
		}
	})

	privateKey, err := os.ReadFile(*privateKeyFilePath)

	if err != nil {
		log.Fatal("Crashed while reading pem from file.")
	}

	jwtToken, err := generateJWT(*applicationID, privateKey)

	if err != nil {
		log.Fatalf("Crashed whilst parsing private key: %s", err)
	}

	token, err := FetchTokenFromAPI(*installationID, jwtToken)

	if err != nil {
		log.Fatalf("Failed to fetch token from API with: %s", err)
	}

	fmt.Print(token)
}

func generateJWT(applicationID string, privateKey []byte) (string, error) {
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix() - 60,
		ExpiresAt: time.Now().Unix() + 120,
		Issuer:    applicationID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	if err != nil {
		return "", err
	}

	jwtToken, err := token.SignedString(signKey)

	if err != nil {
		return "", err
	}

	return jwtToken, err
}

func FetchTokenFromAPI(installationID, jwtToken string) (string, error) {
	bearer := fmt.Sprintf("Bearer %s", jwtToken)

	accessTokenEndpoint := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationID)

	req, _ := http.NewRequest(http.MethodPost, accessTokenEndpoint, nil)

	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return "", err
	}

	if response.StatusCode != 201 {
		return "", errors.New(response.Status)
	}

	defer response.Body.Close()
	var tokenResponse installationTokenResponse
	err = json.NewDecoder(response.Body).Decode(&tokenResponse)

	if err != nil {
		return "", err
	}

	return tokenResponse.Token, nil
}
