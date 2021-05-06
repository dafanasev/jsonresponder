package jsonresponder

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

func TestJSONResponder(t *testing.T) {
	t.Run("response struct", func(t *testing.T) {
		resp := Response{
			HTTPCode: http.StatusAccepted,
			Data: map[string]interface{}{
				"one": 1,
				"two": "second",
			},
			Errors: nil,
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/respond", nil)
		render.Respond(w, r, resp)

		status, ok := r.Context().Value(render.StatusCtxKey).(int)
		assert.True(t, ok)
		assert.Equal(t, http.StatusAccepted, status)
		assert.Equal(t, "{\"data\":{\"one\":1,\"two\":\"second\"}}\n", w.Body.String())
	})

	t.Run("build response", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/respond", nil)
		render.Respond(w, r, BuildResponse(http.StatusContinue, "to be continued..."))

		status, ok := r.Context().Value(render.StatusCtxKey).(int)
		assert.True(t, ok)
		assert.Equal(t, http.StatusContinue, status)
		assert.Equal(t, "{\"data\":\"to be continued...\"}\n", w.Body.String())
	})

	t.Run("single error", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/respond", nil)
		render.Respond(w, r, errors.New("some error"))

		status, ok := r.Context().Value(render.StatusCtxKey).(int)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, "{\"errors\":[{\"description\":\"some error\"}]}\n", w.Body.String())
	})

	t.Run("multiple errors", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/respond", nil)
		render.Respond(w, r, []error{errors.New("some error"), errors.New("second error")})

		status, ok := r.Context().Value(render.StatusCtxKey).(int)
		assert.True(t, ok)
		assert.Equal(t, http.StatusInternalServerError, status)
		assert.Equal(t, "{\"errors\":[{\"description\":\"some error\"},{\"description\":\"second error\"}]}\n", w.Body.String())
	})

	t.Run("map, array, string, int, struct, array of structs", func(t *testing.T) {
		resps := map[string]interface{}{
			"map":    map[string]interface{}{"1": 1, "two": "two"},
			"slice":  []interface{}{1, "two"},
			"string": "the-response",
			"int":    1,
			"struct": struct {
				N int    `json:"n"`
				S string `json:"s"`
			}{N: 1, S: "s"},
			"multiple_structs": []struct {
				N int    `json:"n"`
				S string `json:"s"`
			}{
				{N: 1, S: "s"},
				{N: 2, S: "str"},
			},
		}

		for key, resp := range resps {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/respond", nil)
			render.Respond(w, r, resp)

			status, ok := r.Context().Value(render.StatusCtxKey).(int)
			assert.True(t, ok)
			assert.Equal(t, http.StatusOK, status)
			switch key {
			case "map":
				assert.Equal(t, "{\"data\":{\"1\":1,\"two\":\"two\"}}\n", w.Body.String())
			case "slice":
				assert.Equal(t, "{\"data\":[1,\"two\"]}\n", w.Body.String())
			case "string":
				assert.Equal(t, "{\"data\":\"the-response\"}\n", w.Body.String())
			case "int":
				assert.Equal(t, "{\"data\":1}\n", w.Body.String())
			case "struct":
				assert.Equal(t, "{\"data\":{\"n\":1,\"s\":\"s\"}}\n", w.Body.String())
			case "multiple_structs":
				assert.Equal(t, "{\"data\":[{\"n\":1,\"s\":\"s\"},{\"n\":2,\"s\":\"str\"}]}\n", w.Body.String())
			}
		}
	})
}

