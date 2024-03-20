package fuego

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultOpenAPIHandler(t *testing.T) {
	t.Run("works with DefaultOpenAPIHandler", func(t *testing.T) {
		s := NewServer()

		s.OutputOpenAPISpec()

		require.NotNil(t, s.OpenAPIConfig.UIHandler)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/swagger/index.html", nil)

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), "OpenAPI specification")
		require.Equal(t, "", w.Header().Get("X-Test-Response"))
	})

	t.Run("wrap DefaultOpenAPIHandler behind a middleware", func(t *testing.T) {
		s := NewServer(
			WithOpenAPIConfig(OpenAPIConfig{
				UIHandler: func(specURL string) http.Handler {
					return dummyMiddleware(DefaultOpenAPIHandler(specURL))
				},
			}),
		)
		s.OutputOpenAPISpec()

		require.NotNil(t, s.OpenAPIConfig.UIHandler)

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/swagger/index.html", nil)

		s.Mux.ServeHTTP(w, r)

		require.Equal(t, 200, w.Code)
		require.Contains(t, w.Body.String(), "OpenAPI specification")
		require.Equal(t, "response", w.Header().Get("X-Test-Response"))
	})
}
