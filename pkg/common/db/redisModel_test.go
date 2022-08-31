package db

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetAccountCode(t *testing.T) {
	uid := "test_uid"
	code := 666666
	err := DB.SetAccountCode(uid, code, 100)
	assert.Nil(t, err)
}
func Test_GetAccountCode(t *testing.T) {
	uid := "test_uid"
	code, err := DB.GetAccountCode(uid)
	assert.Nil(t, err)
	fmt.Println(code)
}
