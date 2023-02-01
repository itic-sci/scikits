package scikits

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.
var signature = "self-defined"

func SetSignature(key string) {
	signature = key
}

// 加密
func JwtEncrypt(jwtBody jwt.MapClaims, expTime time.Duration) string {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	jwtBody["exp"] = time.Now().Add(expTime).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtBody)

	// Sign and get the complete encoded token as a string using the secret
	//fmt.Println(signature)
	tokenString, err := token.SignedString([]byte(signature))
	fmt.Println(tokenString, err)
	return tokenString
}

// 解密
func JwtDecrypt(tokenString string) (map[string]interface{}, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(signature), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("Token is expired")
	}
}
