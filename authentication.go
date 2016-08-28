package main

import (
    "github.com/joaodias/hugitoApp/utils"
    "github.com/mitchellh/mapstructure"
    "golang.org/x/oauth2"
    githuboauth "golang.org/x/oauth2/github"
)

// Authentication represents the data exchanged between the client and the server
type Authentication struct {
    Authenticated bool   `json:"authenticated"`
    ClientID      string `json:"clientId"`
    Secret        string `json:"secret"`
    Code          string `json:"code"`
    State         string `json:"state"`
    ReceivedState string `json:"receivedState"`
    AccessToken   string `json:"token"`
}

// authenticate is the main call to perform the authentication process.
func authenticate(client *Client, data interface{}) {
    var authentication Authentication
    err := mapstructure.Decode(data, &authentication)
    if err != nil {
        client.send <- Message{"error", err.Error()}
        return
    }

    go func() {
        if !isStateValid(authentication.State, authentication.ReceivedState) {
            client.send <- Message{"authenticated set", false}
            return
        }

        accessToken, err := getToken(authentication)
        if err != nil {
            client.send <- Message{"error", err.Error()}
            return
        }

        authentication.AccessToken = accessToken
        client.send <- Message{"authenticated set", "true"}
    }()
}

func isStateValid(state, receivedState string) bool {
    return utils.AreStringsEqual(state, receivedState)
}

func getToken(authentication Authentication) (string, error) {
    var (
        oauthConf = &oauth2.Config{
            ClientID:     authentication.ClientID,
            ClientSecret: authentication.Secret,
            Endpoint:     githuboauth.Endpoint,
        }
    )

    token, err := oauthConf.Exchange(oauth2.NoContext, authentication.Code)
    if err != nil {
        return "", err
    }

    return token.AccessToken, nil
}
