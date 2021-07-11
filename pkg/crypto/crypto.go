package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

const encryptionKey ="-JaNcRfUjXn2r5u8x/A?D(G+KbPeSgVk"

func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)

	padLength := aes.BlockSize - len(plainTextBytes)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padLength)}, padLength)
	plainTextBytes = append(plainTextBytes, padText...)

	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]

	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainTextBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(cipherText string) (string, error) {
		block, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		return "", err
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	mod := len(cipherTextBytes) % aes.BlockSize
	if mod != 0 {
		return "", errors.New("blocksize not correct")
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	plainTextBytes := make([]byte, len(cipherTextBytes))
	mode.CryptBlocks(plainTextBytes, cipherTextBytes)

	length := len(plainTextBytes)
	padLength := int(plainTextBytes[length-1])
	plainTextBytes = plainTextBytes[:(length - padLength)]

	return string(plainTextBytes), nil
}
