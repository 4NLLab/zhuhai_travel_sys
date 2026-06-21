package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	hash := "$2a$10$UrVLkzEgC6uk7.wLVYzIy.VDoiK2d5/oc0YkwG3JWQ2E/kI2c/2wq"
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("admin123"))
	fmt.Printf("result: err=%v\n", err)
}
