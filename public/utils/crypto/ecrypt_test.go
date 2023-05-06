package crypto

import (
	"fmt"
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	password := GeneralPwd("admin123")
	fmt.Println(password)
	hash := CheckPwd("admin123", "$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2")
	fmt.Println(hash)

	//定义明文
	data := []byte("15988886666")
	//密钥
	key_aes := []byte("1234567890aaaaaa")
	key_des := []byte("12345678")

	AESEncrypt(data, key_aes)
	AESDecrypt("CqtzrUgSjSNPKHhYVEnkKw==", key_aes)

	DESEncrypt(data, key_des)
	DESDecrypt("HPFOVaFHuj5vs51P/kqmVg==	", key_des)

	fmt.Println(SHA256("15988886666"))

	fmt.Println(MD5("15988886666"))
}
