package jsonresponder

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/hashicorp/go-multierror"
)

func init() {
	render.Respond = JSONResponder
}

type Response struct {
	HTTPCode int         `json:"-"`
	Data     interface{} `json:"data,omitempty"`
	Errors   []Error     `json:"errors,omitempty"`
}

type Error struct {
	Description string `json:"description"`
	Code        int    `json:"code,omitempty"`
	Details     string `json:"details,omitempty"`
}

func JSONResponder(w http.ResponseWriter, r *http.Request, v interface{}) {
	if v == nil {
		render.Status(r, http.StatusNoContent)
		return
	}

	switch val := v.(type) {
	case Response:
		render.Status(r, val.HTTPCode)
		render.JSON(w, r, val)
	case *multierror.Error:
		JSONResponder(w, r, BuildErrorsResponse(http.StatusInternalServerError, val.Errors...))
	case error:
		JSONResponder(w, r, BuildErrorsResponse(http.StatusInternalServerError, val))
	case []error:
		JSONResponder(w, r, BuildErrorsResponse(http.StatusInternalServerError, val...))
	default:
		JSONResponder(w, r, BuildResponse(http.StatusOK, val))
	}

}

func BuildErrorsResponse(httpCode int, errs ...error) Response {
	return BuildResponse(httpCode, nil, errs...)
}

func BuildResponse(httpCode int, data interface{}, errs ...error) Response {
	resp := Response{
		HTTPCode: httpCode,
		Data:     data,
	}

	if len(errs) > 0 {
		resp.Errors = make([]Error, 0, len(errs))
	}
	for _, err := range errs {
		resp.Errors = append(resp.Errors, Error{
			Description: err.Error(),
		})
	}

	return resp
}

