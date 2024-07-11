package util

import (
	"github.com/labstack/echo/v4"
	"strings"
)

func ReplaceServerURL(url string, servers []string, replacement string) string {
	found := ""
	for _, s := range servers {
		if strings.HasPrefix(url, s) {
			found = s
		}
	}
	return strings.Replace(url, found, replacement, 1)
}

func EnrichURL(ctx echo.Context) {
	url := ctx.Request().URL
	url.Scheme = ctx.Scheme()
	url.Host = ctx.Request().Host
}
