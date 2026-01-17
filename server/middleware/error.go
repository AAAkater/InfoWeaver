package middleware

import (
	"encoding/json"
	"net/http"
	"server/models/response"
	"server/utils"

	"github.com/labstack/echo/v4"
)

func CustomHTTPErrorHandler(err error, ctx echo.Context) {

	req := ctx.Request()
	resp := ctx.Response()
	if resp.Committed {
		return
	}

	he, ok := err.(*echo.HTTPError)
	if ok {
		if he.Internal != nil {
			if herr, ok := he.Internal.(*echo.HTTPError); ok {
				he = herr
			}
		}
	} else {
		he = &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	code := he.Code
	message := he.Message

	switch m := he.Message.(type) {
	case string:
		message = echo.Map{"code": response.ERROR, "msg": m, "data": nil}
	case json.Marshaler:
		// do nothing - this type knows how to format itself to JSON
	case error:
		message = echo.Map{"code": response.ERROR, "msg": m.Error(), "data": nil}
	}

	// Send response
	if req.Method == http.MethodHead {
		err = ctx.NoContent(he.Code)
	} else {
		err = ctx.JSON(code, message)
	}
	if err != nil {
		utils.Logger.Error(err)
	}
}
