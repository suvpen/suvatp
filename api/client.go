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

type DidDoc struct {
	Context            []string `json:"@context"`
	Id                 string   `json:"id"`
	AlsoKnownAs        []string `json:"alsoKnownAs"`
	VerificationMethod []struct {
		Id                 string `json:"id"`
		Type               string `json:"type"`
		Controller         string `json:"controller"`
		PublicKeyMultibase string `json:"publicKeyMultibase"`
	} `json:"verificationMethod"`
	Service []struct {
		Id              string `json:"id"`
		Type            string `json:"type"`
		ServiceEndpoint string `json:"serviceEndpoint"`
	} `json:"service"`
}

type Config struct {
	ATProtoEndpoint    string `json:"at_proto_endpoint"`
	PDSEndpoint        string `json:"pds_endpoint"`
	ProfilesCollection string `json:"profiles_collection"`
	PostsCollection    string `json:"posts_collection"`
	RepostsCollection  string `json:"reposts_collection"`
	LikesCollection    string `json:"likes_collection"`
	GraphFollowLexicon string `json:"graph_follow_lexicon"`
	GraphBlockLexicon  string `json:"graph_block_lexicon"`
	LabelerService     string `json:"labeler_service"`
}

type ATPClient struct {
	Config        *Config `json:"config"`
	Client        *xrpc.Client
	PdsClient     *xrpc.Client
	LabelerClient *xrpc.Client
	Did           string
	AppPassword   string
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
	atpClient.PdsClient.Client = nil
	atpClient.LabelerClient.Client = nil

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

	atpClient.PdsClient.Auth = atpClient.Client.Auth
	atpClient.LabelerClient.Auth = atpClient.Client.Auth

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

	return time.Now().Add(time.Minute).Unix() >= jwt.Exp, nil
}

func createSession(did, appPassword, clientAuthFilePath string, config *Config) (*ATPClient, error) {
	atpClient := &ATPClient{
		Config: config,
		Client: &xrpc.Client{
			Client: new(http.Client),
			Host:   config.ATProtoEndpoint,
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

	resultJson, err := json.Marshal(*session.DidDoc)
	if err != nil {
		return nil, err
	}

	var didDoc *DidDoc
	err = json.Unmarshal(resultJson, &didDoc)
	if err != nil {
		return nil, err
	}

	//ATPROTO CLIENT
	atpClient.Client.Auth = &xrpc.AuthInfo{
		AccessJwt:  session.AccessJwt,
		RefreshJwt: session.RefreshJwt,
		Handle:     session.Handle,
		Did:        session.Did,
	}

	//PDS CLIENT
	atpClient.Config.PDSEndpoint = didDoc.Service[0].ServiceEndpoint
	atpClient.PdsClient = &xrpc.Client{
		Client: new(http.Client),
		Host:   didDoc.Service[0].ServiceEndpoint,
	}
	atpClient.PdsClient.Auth = atpClient.Client.Auth

	seeds := make(map[string]string)
	seeds["Atproto-Proxy"] = "did:web:api.bsky.chat#bsky_chat"
	atpClient.PdsClient.Headers = seeds

	//LABELER CLIENT
	atpClient.LabelerClient = &xrpc.Client{Client: new(http.Client)}
	if len(didDoc.Service) > 1 {
		atpClient.LabelerClient.Host = didDoc.Service[0].ServiceEndpoint
		atpClient.LabelerClient.Auth = atpClient.Client.Auth

		seeds = make(map[string]string)
		seeds["Atproto-Proxy"] = "did:plc:yojwcfgpkxq35sv5wioglqad#atproto_labeler"
		atpClient.LabelerClient.Headers = seeds
	}

	err = writeAuthFile(clientAuthFilePath, *atpClient)
	if err != nil {
		return nil, err
	}

	return atpClient, nil
}

func getClientAuthFile(atpEndpoint, did string) (string, error) {
	atpName := strings.Replace(atpEndpoint, "https://", "", 1)
	didFileName := strings.Replace(did, "did:plc:", "", 1)
	clientAuthFilePath := fmt.Sprintf(ATPClientAuthJsonFile, atpName, didFileName)

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
			ATProtoEndpoint:    DefaultATProtoEndpoint,
			ProfilesCollection: DefaultProfilesCollection,
			PostsCollection:    DefaultPostsCollection,
			RepostsCollection:  DefaultRepostsCollection,
			LikesCollection:    DefaultLikeCollection,
			GraphFollowLexicon: DefaultGraphFollowLexicon,
			GraphBlockLexicon:  DefaultGraphBlockLexicon,
			LabelerService:     DefaultLabelerService,
		}
	}

	clientAuthFilePath, err := getClientAuthFile(config.ATProtoEndpoint, did)
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

		if appPassword != atpClient.AppPassword || atpClient.Config != config {
			atpClient.AppPassword = appPassword
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
		atpClient.PdsClient.Client = new(http.Client)
		atpClient.LabelerClient.Client = new(http.Client)

		if jwtIsExpired {
			atpClient, err = refreshSession(atpClient, clientAuthFilePath)
			if err != nil {
				return nil, err
			}
		}
	}

	return atpClient, nil
}
