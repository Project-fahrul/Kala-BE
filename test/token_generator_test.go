package test

import (
	"errors"
	"fmt"
	"kala/service"
	"kala/util"
	"testing"
)

func Test_generatorToken(t *testing.T) {
	token, err := util.TokenGenerator(9, "fahrul@gmail.com")

	if err != nil {
		t.Log(err.Error())
		t.Error()
	}

	t.Logf("\nToken gen: %v\n", token)

	token, err = util.Decrypt(token)
	if err != nil {
		t.Log(err.Error())
		t.Error()
	}

	t.Logf("\nToken dec: %v\n", token)
}

func validateUserWithToken(token, email, keyword string, withValidateToken bool) error {

	_token, err := service.Redis_New().Get(fmt.Sprintf("%s:%s", email, keyword))

	if err != nil {
		return err
	}

	if _token != token {
		return errors.New("Unknow token")
	}

	if withValidateToken {
		err = util.ValidateToken(_token, email)
		if err != nil {
			return err
		}
	}

	return nil
}
