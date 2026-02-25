package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
)

func CustomHTTPErrorHandler() echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		if r, _ := echo.UnwrapResponse(c.Response()); r != nil && r.Committed {
			return
		}

		code := http.StatusInternalServerError
		var sc echo.HTTPStatusCoder
		if errors.As(err, &sc) {
			if tmp := sc.StatusCode(); tmp != 0 {
				code = tmp
			}
		}

		var result any
		switch m := sc.(type) {
		case json.Marshaler: // this type knows how to format itself to JSON
			result = m
		case *echo.HTTPError:
			sText := m.Message
			if sText == "" {
				sText = http.StatusText(code)
			}
			msg := map[string]any{
				"code": 1,
				"msg":  sText,
				"data": nil,
			}
			result = msg
		default:
			msg := map[string]any{
				"code": 1,
				"msg":  http.StatusText(code),
				"data": nil,
			}
			result = msg
		}

		var cErr error
		if c.Request().Method == http.MethodHead { // Issue #608
			cErr = c.NoContent(code)
		} else {
			cErr = c.JSON(code, result)
		}
		if cErr != nil {
			c.Logger().Error("echo default error handler failed to send error to client", "error", cErr) // truly rare case. ala client already disconnected
		}
	}
}
