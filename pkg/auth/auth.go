package auth

import (
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/dghubble/oauth1/twitter"
	"github.com/niceb5y/heart-extractor/pkg/crypto"
	"log"
)


const EncryptedConsumerKey = "SGuVNhwL9dUL6BfigQzMF4J8vwkwfNrsgfE3zKK8h94T0UJGZQV2jICgLmkGIH4O"
const EncryptedConsumerSecret = "ArMhqPeuAXk3ad+n3dEOW0gLwpcOLaSpD02H80HY2ahhEHXOUMUdlsEWjW7wdUfvHIKKFz83xyncrAj/M6j7qmCKzGi0gcgJYKsqfI12xuQ="

func CreateAuthToken() *oauth1.Token {
	consumerKey, _ := crypto.Decrypt(EncryptedConsumerKey)
	consumerSecret, _ := crypto.Decrypt(EncryptedConsumerSecret)

	config := oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    "oob",
		Endpoint:       twitter.AuthorizeEndpoint,
	}
	requestToken, err := login(config)
	if err != nil {
		log.Fatalln(err)
	}

	accessToken, err := getAuthTokenWithPIN(config, requestToken)
	if err != nil {
		log.Fatalln(err)
	}

	return accessToken
}

func login(config oauth1.Config) (requestToken string, err error) {
	requestToken, _, err = config.RequestToken()
	if err != nil {
		return "", err
	}

	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		return "", err
	}

	fmt.Println("Access token needed")
	fmt.Println("Open this URL for authentication:\n", authorizationURL.String())
	return requestToken, err
}

func getAuthTokenWithPIN(config oauth1.Config, requestToken string) (*oauth1.Token, error) {
	print("Enter PIN: ")
	var pin string
	_, err := fmt.Scanf("%s", &pin)
	if err != nil {
		return nil, err
	}

	const tokenRequestSecret = "Yq3t6w9z$C&F)J@NcQfTjWnZr4u7x!A%"
	accessToken, accessSecret, err := config.AccessToken(requestToken, tokenRequestSecret, pin)
	if err != nil {
		return nil, err
	}

	return oauth1.NewToken(accessToken, accessSecret), err
}
