package main

import (
	"bytes"
	"github.com/gorilla/rpc/json"
	"github.com/ilyail3/go-json-rpc-example/shared"
	"log"
	"net/http"
)

func main() {
	url := "http://localhost:1234/rpc"

	args := &shared.Args{
		A: 2,
		B: 4,
	}

	message, err := json.EncodeClientRequest("Arith.Multiply", args)
	if err != nil {
		log.Fatalf("%s", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		log.Fatalf("%s", err)
	}

	tp := shared.NewTokenProvider([]byte("AllYourBase"))
	err, ss := tp.AccountToken("account")

	if err != nil {
		log.Fatalf("error creating token: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Jwt-Auth", ss)
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error in sending request to %s. %s", url, err)
	}

	defer resp.Body.Close()

	var result shared.Result
	err = json.DecodeClientResponse(resp.Body, &result)

	if err != nil {
		log.Fatalf("Couldn't decode response. %s", err)
	}

	log.Printf("%d*%d=%d\n", args.A, args.B, result)
}