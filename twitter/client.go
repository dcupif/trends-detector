package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

const (
	baseURL           string = "https://api.twitter.com"
	auth                     = "/oauth2/token"
	streamFilter             = "/labs/1/tweets/stream/filter"
	streamFilterRules        = streamFilter + "/rules"
)

type Client struct {
	*http.Client
}

func NewClient(key, secret string) (Client, error) {
	conf := oauth2.Config{
		ClientID:     key,
		ClientSecret: secret,
	}
	token, err := bearerToken(key, secret)
	if err != nil {
		panic(err)
	}
	return Client{Client: conf.Client(context.Background(), token)}, nil
}

func bearerToken(key, secret string) (*oauth2.Token, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+auth,
		strings.NewReader("grant_type=client_credentials"),
	)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(key, secret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	var client http.Client
	res, err := client.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, err
	}

	var token oauth2.Token
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

type Rule struct {
	Value string `json:"value"`
	Tag   string `json:"tag, omitempty"`
	ID    string `json:"id, omitempty"`
}

type RulePayload struct {
	Add []Rule `json:"add"`
}

type RuleResponse struct {
	Data []Rule `json:"data"`
}

func (c *Client) StreamRules() (string, error) {
	resp, err := c.Get(baseURL + streamFilterRules)
	if err != nil {
		return "", err
	}

	rules, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(rules), nil
}

func (c *Client) AddStreamRule(rules []Rule, dryRun bool) ([]Rule, error) {
	var payload bytes.Buffer
	enc := json.NewEncoder(&payload)
	err := enc.Encode(RulePayload{Add: rules})
	if err != nil {
		return nil, fmt.Errorf("failed to serialize rules %+v: %v", rules, err)
	}

	baseUrl, err := url.Parse(baseURL + streamFilterRules)
	if err != nil {
		return nil, fmt.Errorf("malformed URL: %v", err)
	}

	if dryRun {
		params := url.Values{}
		params.Add("dry_run", "true")
		baseUrl.RawQuery = params.Encode() // Escape Query Parameters
	}

	resp, err := c.Post(
		baseUrl.String(),
		"application/json",
		&payload,
	)
	if err != nil {
		return nil, err
	}

	var newRulesResp RuleResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&newRulesResp)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize rules: %v", err)
	}

	return newRulesResp.Data, nil
}
