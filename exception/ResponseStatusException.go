package exception

import (
	"net/http"
)

type ResponseStatusException struct {
	Error error
	Code  int
}

func ResponseStatusError_New(err error) {
	if err != nil {
		panic(ResponseStatusException{
			Error: err,
			Code:  http.StatusBadRequest,
		})
	}
}
