package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/itic-sci/scikits"
	"time"
)

func test() {
	body := jwt.MapClaims{
		"user": "username",
		"iat":  time.Now().Hour(),
	}

	expTime := time.Second * time.Duration(1)

	ts := scikits.JwtEncrypt(body, expTime)
	//time.Sleep(time.Second * time.Duration(3))

	claims, err := scikits.JwtDecrypt(ts)

	fmt.Println(claims, err)
}

func main() {
	//test()

	//time.Sleep(time.Second * time.Duration(2))

	//scikits.SetSignature("self-define")
	//
	//body := jwt.MapClaims{
	//	"user": "username",
	//	"iat":  time.Now().Hour(),
	//}
	//
	//fmt.Println(body)
	//
	//expTime := time.Second * time.Duration(1)
	//
	//ts := scikits.JwtEncrypt(body, expTime)
	//fmt.Println(ts)
	//

	scikits.SetSignature("wqeqe123123")
	tss := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE2OTExNDEyNjYsImV4cCI6MTY5MTE0MTQ2Nn0.gmgHNAodSch25Hp256Pez_miFq0sZop1QqB3k5Stt4"
	claims, err := scikits.JwtDecrypt(tss)

	fmt.Println(claims, err)
}
