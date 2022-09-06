package model

import "github.com/gin-gonic/gin"

func private_HTTPResponse_Generate(status bool, data interface{}, msg ...string) gin.H {
	res := gin.H{}
	res["status"] = status

	if status && data != nil {
		res["data"] = data
	} else if !status && len(msg) > 0 {
		res["message"] = msg[0]
	}

	return res
}

func HTTPResponse_Message(msg string) gin.H {
	return private_HTTPResponse_Generate(false, nil, msg)
}

func HTTPResponse_Data(data interface{}) gin.H {
	return private_HTTPResponse_Generate(true, data)
}
