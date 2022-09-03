package model

import "github.com/gin-gonic/gin"

func private_HTTPResponse_Generate(status bool, data interface{}, msg ...string) interface{} {
	res := gin.H{}
	res["status"] = status

	if status && data != nil {
		res["data"] = data
	} else if !status && len(msg) > 0 {
		res["message"] = msg[0]
	}

	return res
}

func HTTPResponse_Message(msg string) interface{} {
	return private_HTTPResponse_Generate(false, nil, msg)
}

func HTTPResponse_Data(data interface{}) interface{} {
	return private_HTTPResponse_Generate(true, data)
}
