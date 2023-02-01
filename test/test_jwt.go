package main

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/xulehexuwei/scikits"
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
	test()

	//time.Sleep(time.Second * time.Duration(2))
	//
	scikits.SetSignature("self-define")

	body := jwt.MapClaims{
		"user": "username",
		"iat":  time.Now().Hour(),
	}

	expTime := time.Second * time.Duration(1)

	ts := scikits.JwtEncrypt(body, expTime)
	//time.Sleep(time.Second * time.Duration(3))

	scikits.SetSignature("self-defined")

	claims, err := scikits.JwtDecrypt(ts)

	fmt.Println(claims, err)
}
