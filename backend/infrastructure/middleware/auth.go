package middleware

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"master-control/model"
)

// TokenQueryParam es el parámetro de query aceptado para recursos que el navegador
// pide sin cabeceras (p. ej. <img src> de los iconos).
const TokenQueryParam = "token"

// Auth exige el token en `Authorization: Bearer <token>` o en `?token=`.
// Las rutas de `public` (p. ej. el health check) y los preflight OPTIONS pasan sin token.
func Auth(token string, public ...string) gin.HandlerFunc {
	expected := []byte(token)
	open := make(map[string]struct{}, len(public))
	for _, p := range public {
		open[p] = struct{}{}
	}

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		if _, ok := open[c.FullPath()]; ok {
			c.Next()
			return
		}

		given := []byte(extractToken(c))
		if len(given) == 0 || subtle.ConstantTimeCompare(given, expected) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Error{Error: "token inválido o ausente"})
			return
		}

		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	if header := c.GetHeader("Authorization"); header != "" {
		if after, ok := strings.CutPrefix(header, "Bearer "); ok {
			return strings.TrimSpace(after)
		}
	}

	return c.Query(TokenQueryParam)
}

// Logger registra las peticiones sin la query string, para no volcar el token en el log.
// `skip` son rutas de alta frecuencia (p. ej. el scroll táctil) que no se registran.
func Logger(skip ...string) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: skip,
		Formatter: func(p gin.LogFormatterParams) string {
			path := p.Path
			if i := strings.IndexByte(path, '?'); i >= 0 {
				path = path[:i]
			}
			p.Path = path
			return defaultFormatter(p)
		},
	})
}

func defaultFormatter(p gin.LogFormatterParams) string {
	line := fmt.Sprintf("[GIN] %s | %3d | %13v | %15s | %-7s %s",
		p.TimeStamp.Format("2006/01/02 - 15:04:05"), p.StatusCode, p.Latency, p.ClientIP, p.Method, p.Path)
	if p.ErrorMessage != "" {
		line += " | " + p.ErrorMessage
	}
	return line + "\n"
}
