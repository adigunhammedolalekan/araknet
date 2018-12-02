package araknet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
)

type Encryptor struct {
	key string
}

func NewEncryptor(key string) *Encryptor {

	return &Encryptor{
		key: key,
	}
}

func (e *Encryptor) Encrypt(data interface{}) ([]byte, error) {

	key := createKeyHash(e.key)
	block, _ := aes.NewCipher([]byte(key))

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	d, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	cipherText := gcm.Seal(nonce, nonce, d, nil)
	return cipherText, nil
}

func (e *Encryptor) Decrypt(data []byte) ([]byte, error) {

	key := createKeyHash(os.Getenv("CARD_ENC_KEY"))
	block, _ := aes.NewCipher([]byte(key))

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	text, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	return text, nil
}

func createKeyHash(key string) string {

	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString([]byte(hasher.Sum(nil)))
}
