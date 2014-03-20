package webapi

import (
	"net/http"
)

type WebApiRequest struct {
	Http			*http.Request
	Params			map[string]string
}