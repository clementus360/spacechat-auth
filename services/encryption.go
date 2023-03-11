package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)



func GenerateSecret() (string,error) {
	secret := make([]byte,20)

	_,err := rand.Read(secret)

	if err!=nil {
		fmt.Println("Failed to generate secret")
		return "",err
	}

	secretKey := base32.StdEncoding.EncodeToString(secret)

	return secretKey[:16],nil
}

func Encrypt(s string, key string) (string,error) {
	// Convert string to byte array
	plainText := []byte(s)

	// Create new AES cypher using the secret key
	block,err := aes.NewCipher([]byte(key))
	if err!=nil{
		return "",err
	}

	// Apply AES cypher to the plaintext
	cypherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cypherText[:aes.BlockSize]
	if _,err := io.ReadFull(rand.Reader,iv); err!=nil {
		return "",err
	}
	stream := cipher.NewCFBEncrypter(block,iv)
	stream.XORKeyStream(cypherText[aes.BlockSize:], plainText)

	return base64.URLEncoding.EncodeToString(cypherText),nil
}

func Decrypt(s string, key string) (string,error) {
	// Convert string cyphertext to bytes
	cypherText, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	// Create new AES cypher using the secret key
	block,err := aes.NewCipher([]byte(key))
	if err!=nil{
		return "",err
	}

	if len(cypherText) < aes.BlockSize {
		return "",errors.New("cyphertext is too short")
	}

	iv := cypherText[:aes.BlockSize]
	cypherText = cypherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	plainText := make([]byte,len(cypherText))
	stream.XORKeyStream(plainText,cypherText)

	return string(plainText),nil
}

func Hash(s string) string {
	hash:= sha256.Sum256([]byte(s))
	return hex.EncodeToString(hash[:])
}
