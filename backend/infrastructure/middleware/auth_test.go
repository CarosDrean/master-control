package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func router() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	g := r.Group("/mc")
	g.Use(Auth("s3cret", "/mc/health"))
	ok := func(c *gin.Context) { c.Status(http.StatusOK) }
	g.GET("/health", ok)
	g.GET("/actions", ok)
	g.OPTIONS("/actions", ok)
	return r
}

func TestAuth(t *testing.T) {
	cases := []struct {
		name   string
		method string
		path   string
		header string
		want   int
	}{
		{"health público", http.MethodGet, "/mc/health", "", http.StatusOK},
		{"sin token", http.MethodGet, "/mc/actions", "", http.StatusUnauthorized},
		{"token incorrecto", http.MethodGet, "/mc/actions", "Bearer nope", http.StatusUnauthorized},
		{"esquema incorrecto", http.MethodGet, "/mc/actions", "Basic s3cret", http.StatusUnauthorized},
		{"bearer válido", http.MethodGet, "/mc/actions", "Bearer s3cret", http.StatusOK},
		{"query válida", http.MethodGet, "/mc/actions?token=s3cret", "", http.StatusOK},
		{"query inválida", http.MethodGet, "/mc/actions?token=x", "", http.StatusUnauthorized},
		{"preflight", http.MethodOptions, "/mc/actions", "", http.StatusOK},
	}
	r := router()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			if w.Code != tc.want {
				t.Errorf("%s %s = %d, quiero %d", tc.method, tc.path, w.Code, tc.want)
			}
		})
	}
}

func TestLoggerOmiteQuery(t *testing.T) {
	var buf bytes.Buffer
	prev := gin.DefaultWriter
	gin.DefaultWriter = &buf
	defer func() { gin.DefaultWriter = prev }()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Logger())
	r.GET("/icon", func(c *gin.Context) { c.Status(http.StatusOK) })
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/icon?token=s3cret&path=x", nil))

	line := buf.String()
	if !strings.Contains(line, "/icon") || strings.Contains(line, "s3cret") {
		t.Errorf("el log debe tener la ruta sin la query: %q", line)
	}
}
