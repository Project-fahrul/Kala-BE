package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"kala/service"
	"strconv"
	"time"
)

const (
	TOKEN_CHANGE_PASSWORD = 1

	keyTokenChangePasswor = "CHANGE_PASSWORD"
	secret                = "abc&1*~#^2^#s0^=)^^7%b34"
	TOKENEXPIRED          = 5
)

func TokenGenerator(tokenType int, email string) (res string, e error) {
	defer func() {
		err := recover()
		if err != nil {
			res = ""
			e = errors.New(err.(string))
		}
	}()
	expired := time.Now().Add(time.Duration(TOKENEXPIRED) * time.Minute).Unix()
	expiredStr := strconv.FormatInt(expired, 10)
	data := map[string]string{
		"email": email,
		"expr":  expiredStr,
	}

	if tokenType == TOKEN_CHANGE_PASSWORD {
		data["key"] = keyTokenChangePasswor
	} else {
		panic("key token invalid")
	}

	js, err := json.Marshal(data)

	if err != nil {
		panic(err.Error())
	}

	enc, err := encrypt(string(js))
	res = enc
	e = err
	return
}

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func ValidateToken(token, email string) error {
	var data map[string]string
	err := json.Unmarshal([]byte(token), &data)

	if err != nil {
		return err
	}

	if data["email"] != email {
		return errors.New("email not match")
	}

	expiredInt, err := strconv.ParseInt(data["expr"], 10, 64)
	unix := time.Unix(expiredInt, 0).UTC().Unix()
	if err != nil || unix > time.Now().UTC().Unix() {
		return errors.New("Token expired")
	}

	return nil
}

func ValidateWithToken(token, email, keyword string, withValidateToken bool) error {

	_token, err := service.Redis_New().Get(fmt.Sprintf("%s:%s", email, keyword))

	if err != nil {
		return err
	}

	if _token != token {
		return errors.New("Unknow token")
	}

	if withValidateToken {
		err = ValidateToken(_token, email)
		if err != nil {
			return err
		}
	}

	return nil
}
