package scikits

import (
	"crypto/sha1"
	"fmt"
	"io"
	"math/rand"
)

func GetRandSalt(length int) string {
	salt := ""
	codeChars := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	for i := 0; i < length; i++ {
		index := rand.Intn(len(codeChars))
		salt = salt + codeChars[index]
	}
	return salt
}

func GetSHA1(saltPassword string) string {
	t := sha1.New()
	io.WriteString(t, saltPassword)
	encodeString := fmt.Sprintf("%x", t.Sum(nil))
	return encodeString
}

// 密码加密, sha1(8位随机盐+明文密码）
func PasswordCrypt(password string) (string, string) {
	// 生成8位随机盐
	salt := GetRandSalt(8)

	// 对字符串进行SHA1哈希
	encodePw := GetSHA1(salt + password)
	return encodePw, salt
}

// 密码验证是否正确
func PasswordVerify(saltPassword, encodePW string) bool {
	sha1String := GetSHA1(saltPassword)
	if sha1String != encodePW {
		return false
	} else {
		return true
	}
}
