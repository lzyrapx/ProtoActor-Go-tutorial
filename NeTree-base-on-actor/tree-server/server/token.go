package server

import (
	"log"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// 创建 length 位 token
func CreateToken(length int) string {
	newRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	tokenBytes := make([]byte, length)
	for i := range tokenBytes {
		tokenBytes[i] = charset[newRand.Intn(len(charset))]
	}
	return string(tokenBytes)
}


func isTokenValid(token string, item TreeItem) bool {
	return token == item.token
}

// 检验 token 合法性
func CheckIDAndToken2(ok bool, id int32, token string, item TreeItem) bool {
	if ok {
		if !isTokenValid(token, item) {
			log.Printf("Tree with id: %v does not correspond to token: %v\n", id, token)
			log.Printf("ID: %v, Token: %v\n", item.id, item.token)
			return false
		}
		return true
	}
	log.Printf("Tree with id: %v does not exist\n", id)
	return false
}
