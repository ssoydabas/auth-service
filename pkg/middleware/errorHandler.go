package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/ssoydabas/auth-service/pkg/errors"
)

func ErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		if appErr, ok := err.(*errors.AppError); ok {
			return c.JSON(appErr.Code, appErr)
		}

		internalErr := errors.InternalError(err)
		return c.JSON(internalErr.Code, internalErr)
	}
}
