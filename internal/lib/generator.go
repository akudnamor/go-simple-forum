package lib

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func GenerateUserId() (string, error) {
	const op = "generator.GenerateUserId"

	buf := make([]byte, 5)
	_, err := rand.Read(buf)

	if err != nil {
		log.Println(err)
		return "", err
	}

	userId := hex.EncodeToString(buf)

	return userId, nil
}
