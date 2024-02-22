package jwt

import (
	"fmt"
	"testing"
)

func TestGenWithParse(t *testing.T) {
	token, err := GenToken(1, 1)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(token)
	claims, err := ParseToken(token)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(claims)
}
