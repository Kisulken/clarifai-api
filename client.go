package clarifai_api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"io/ioutil"
	"strings"
	"net/url"
	"encoding/base64"
	"sync/atomic"
	"errors"
	"bytes"
	"time"
)

const (
	rootURLv2 = "https://api.clarifai.com/v2/"
	rootURLv1 = "https://api.clarifai.com/v1/"

	GeneralModelID = "aaa03c23b3724a16a56b629203edc62c"
	CelebrityModelID = "e466caa0619f444ab97497640cefc4dc"
	ApparelModelID = "e0be3b9d6a454f0493ac3a30784001ff"
	NsfwModelID = "e9576d86d2004ed1a38ba0cf39ecb4b1" // not safe for work model
	FoodModelID = "bd367be194cf45149e75f01d59f77ba7"
)

type ClientConfig struct {
	id             string
	secret         string
	basicAuthToken string
}

type Client struct {
	config	    atomic.Value
	accessToken atomic.Value
	http	    *http.Client
}

func NewClient(clientID, clientSecret string) *Client {
	basicAuthToken := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	config := ClientConfig{clientID, clientSecret, basicAuthToken}
	client := new(Client)
	client.config.Store(config)
	transport := &http.Transport{
		MaxIdleConns: 10,
		IdleConnTimeout: 20 * time.Second,
		DisableCompression: true,
	}
	client.http = &http.Client{Transport: transport}
	return client
}

func (client *Client) SetConfig(id, secret string) {
	config := client.config.Load().(ClientConfig)
	config.id = id
	config.secret = secret
	config.basicAuthToken = base64.StdEncoding.EncodeToString([]byte(id + ":" + secret))
	client.config.Store(config)
}

func (client *Client) SetTransport(transport *http.Transport) {
	client.http.Transport = transport
}

func (client *Client) CustomRequest(root, endpoint, method string, body []byte) ([]byte, error) {
	req, err := http.NewRequest(method, root + endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	// check access token and if it does not exist, then request a new one
	var token TokenResponse
	tokenInterface := client.accessToken.Load()
	if tokenInterface == nil {
		err = client.requestAccessToken()
		if err != nil {
			return nil, err
		}
		token = client.accessToken.Load().(TokenResponse)
	} else {
		token = tokenInterface.(TokenResponse)
	}

	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	req.Header.Set("Authorization", "Bearer " + token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.http.Do(req)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 200, 201:
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		return body, err
	case 401:
		return nil, errors.New("TOKEN_INVALID")
	case 400:
		return nil, errors.New("ALL_ERROR")
	case 404:
		return nil, errors.New("404 NOT FOUND")
	case 500:
		return nil, errors.New("CLARIFAI_ERROR")
	default:
		return nil, errors.New("UNEXPECTED_STATUS_CODE " + strconv.Itoa(res.StatusCode))
	}
}

func (client *Client) requestAccessToken() error {
	configInterface := client.config.Load()
	if configInterface == nil {
		return errors.New("Configuration was not set up")
	}
	config := configInterface.(ClientConfig)

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", config.id)
	form.Set("client_secret", config.secret)
	formData := strings.NewReader(form.Encode())

	req, err := http.NewRequest("POST", rootURLv2 + "token", formData)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Basic " + config.basicAuthToken)
	req.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	token := TokenResponse{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return err
	}

	client.accessToken.Store(token)

	return nil
}