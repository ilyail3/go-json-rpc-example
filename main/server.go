package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"fmt"
	"errors"
)

type Args struct {
	A, B int
}

type Arith struct{
	secret []byte
}

func (t *Arith) authenticate(r *http.Request) error {
	tokenString := r.Header.Get("Jwt-Auth")

	if tokenString == "" {
		return errors.New("header Jwt-Auth is blank or missing")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return t.secret, nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		log.Println(claims)
		return nil
	} else {
		return err
	}
}

type Result int

func (t *Arith) Multiply(r *http.Request, args *Args, result *Result) error {
	err := t.authenticate(r)
	if err != nil {
		return err
	}

	log.Printf("Multiplying %d with %d\n", args.A, args.B)
	*result = Result(args.A * args.B)
	return nil
}

func main() {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	arith := &Arith{ secret:[]byte("AllYourBase") }
	s.RegisterService(arith, "")
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	http.ListenAndServe(":1234", r)
}