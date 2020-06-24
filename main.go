package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dcupif/trendetector/twitter"
)

const secretFilename = ".keys.json"

type APICredentials struct {
	Key    string `json:"consumer_key"`
	Secret string `json:"consumer_secret"`
}

func main() {
	credentials, err := apiCredentials()
	if err != nil {
		panic(err)
	}

	client, err := twitter.NewClient(credentials.Key, credentials.Secret)
	if err != nil {
		panic(err)
	}

	rules, err := client.StreamRules()
	if err != nil {
		panic(err)
	}

	fmt.Println("Rules: " + rules)

	// rule := twitter.Rule{
	// 	Value: "cat has:media",
	// 	Tag:   "cats with media",
	// }
	//
	// newRules, err := client.AddStreamRule([]twitter.Rule{rule}, true)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Printf("Rules to be added: %+v\n", newRules)
	//
	// rules, err = client.StreamRules()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// fmt.Println("New rules: " + rules)
}

func apiCredentials() (APICredentials, error) {
	f, err := os.Open(secretFilename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var credentials struct {
		Key    string `json:"consumer_key"`
		Secret string `json:"consumer_secret"`
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&credentials)

	return credentials, err
}
