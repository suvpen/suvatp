package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	Pds                string `json:"pds"`
	ProfilesCollection string `json:"profiles_collection"`
	PostsCollection    string `json:"posts_collection"`
	RepostsCollection  string `json:"reposts_collection"`
	LikesCollection    string `json:"likes_collection"`
	GraphFollowLexicon string `json:"graph_follow_lexicon"`
	GraphBlockLexicon  string `json:"graph_block_lexicon"`
}

type ATPClient struct {
	Config      *Config `json:"config"`
	Client      *xrpc.Client
	Did         string
	AppPassword string
}

type Jwt struct {
	Scope string `json:"scope"`
	Sub   string `json:"sub"`
	Iat   int    `json:"iat"`
	Exp   int64  `json:"exp"`
	Aud   string `json:"aud"`
}

func writeAuthFile(clientAuthFilePath string, atpClient ATPClient) error {
	atpClient.Client.Client = nil

	clientAuthJson, err := json.Marshal(atpClient)
	if err != nil {
		return fmt.Errorf(
			"error marshalling %s: %w", clientAuthFilePath, err)
	}

	if err = os.WriteFile(clientAuthFilePath, clientAuthJson, 0666); err != nil {
		return fmt.Errorf("error writing %s: %w", clientAuthFilePath, err)
	}

	return nil
}

func refreshSession(atpClient *ATPClient, clientAuthFilePath string) (*ATPClient, error) {
	atpClient.Client.Auth.AccessJwt = atpClient.Client.Auth.RefreshJwt

	refresh, err := atproto.ServerRefreshSession(context.TODO(), atpClient.Client)
	if err != nil {
		return nil, err
	}

	atpClient.Client.Auth.Did = refresh.Did
	atpClient.Client.Auth.AccessJwt = refresh.AccessJwt
	atpClient.Client.Auth.RefreshJwt = refresh.RefreshJwt

	err = writeAuthFile(clientAuthFilePath, *atpClient)
	if err != nil {
		return nil, err
	}

	return atpClient, nil
}

func getJWTExpiration(atpClient *ATPClient, clientAuthFilePath string) (bool, error) {
	parts := strings.Split(atpClient.Client.Auth.AccessJwt, ".")
	payloadJson, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false, fmt.Errorf("error decoding %s jwtJson: %w", clientAuthFilePath, err)
	}

	var jwtJson bytes.Buffer
	if err = json.Indent(&jwtJson, payloadJson, "", "  "); err != nil {
		return false, err
	}

	var jwt *Jwt
	err = json.Unmarshal([]byte(jwtJson.String()), &jwt)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling JWT of %s: %w", clientAuthFilePath, err)
	}

	return time.Now().Unix() >= jwt.Exp, nil
}

func createSession(did, appPassword, clientAuthFilePath string, config *Config) (*ATPClient, error) {
	atpClient := &ATPClient{
		Client: &xrpc.Client{
			Client: new(http.Client),
			Host:   config.Pds,
		},
		Did:         did,
		AppPassword: appPassword,
	}

	sessionInput := &atproto.ServerCreateSession_Input{
		Identifier: did,
		Password:   appPassword,
	}

	session, err := atproto.ServerCreateSession(context.TODO(), atpClient.Client, sessionInput)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}

	atpClient.Config = config
	atpClient.Client.Auth = &xrpc.AuthInfo{
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
		Handle:     session.Handle,
		Did:        session.Did,
	}

	err = writeAuthFile(clientAuthFilePath, *atpClient)
	if err != nil {
		return nil, err
	}

	return atpClient, nil
}

func getClientAuthFile(pds, did string) (string, error) {
	pdsName := strings.Replace(pds, "https://", "", 1)
	didFileName := strings.Replace(did, "did:plc:", "", 1)
	clientAuthFilePath := fmt.Sprintf(ATPClientAuthJsonFile, pdsName, didFileName)

	_ = os.Mkdir(ATPDir, os.ModePerm)

	if _, err := os.Stat(clientAuthFilePath); err != nil {
		_, err = os.Create(clientAuthFilePath)
		if err != nil {
			return "", fmt.Errorf("error creating %s: %w", clientAuthFilePath, err)
		}
	}

	return clientAuthFilePath, nil
}

func Client(did, appPassword string, config *Config) (*ATPClient, error) {
	var atpClient *ATPClient

	if config == nil {
		config = &Config{
			Pds:                DefaultPDS,
			ProfilesCollection: DefaultProfilesCollection,
			PostsCollection:    DefaultPostsCollection,
			RepostsCollection:  DefaultRepostsCollection,
			LikesCollection:    DefaultLikeCollection,
			GraphFollowLexicon: DefaultGraphFollowLexicon,
			GraphBlockLexicon:  DefaultGraphBlockLexicon,
		}
	}

	clientAuthFilePath, err := getClientAuthFile(config.Pds, did)
	if err != nil {
		return nil, err
	}

	fileContent, err := os.ReadFile(clientAuthFilePath)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", clientAuthFilePath, err)
	}

	if string(fileContent) == "" {
		atpClient, err = createSession(did, appPassword, clientAuthFilePath, config)
		if err != nil {
			return nil, err
		}
	} else {
		if err = json.Unmarshal(fileContent, &atpClient); err != nil {
			return nil, fmt.Errorf("error unmarshalling %s: %w", clientAuthFilePath, err)
		}

		if atpClient.Config != config {
			atpClient.Config = config

			err = writeAuthFile(clientAuthFilePath, *atpClient)
			if err != nil {
				return nil, err
			}
		}

		jwtIsExpired, err := getJWTExpiration(atpClient, clientAuthFilePath)
		if err != nil {
			return nil, err
		}

		atpClient.Client.Client = new(http.Client)

		if jwtIsExpired {
			atpClient, err = refreshSession(atpClient, clientAuthFilePath)
			if err != nil {
				return nil, err
			}
		}
	}

	return atpClient, nil
}
