package handler

import (
	"net/http"
	"strconv"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/internal/service"
	"github.com/ssoydabas/auth-service/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AccountHandler interface {
	AddRoutes(e *echo.Group)

	CreateAccount(c echo.Context) error
	AuthenticateAccount(c echo.Context) error
	GetAccountByID(c echo.Context) error
	GetAccountByToken(c echo.Context) error
}

type accountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler(accountService service.AccountService) AccountHandler {
	return &accountHandler{
		accountService: accountService,
	}
}

func (h *accountHandler) AddRoutes(e *echo.Group) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.POST("/accounts", h.CreateAccount)
	e.GET("/accounts/:id", h.GetAccountByID)
	e.POST("/accounts/authenticate", h.AuthenticateAccount)
	e.GET("/accounts/me", h.GetAccountByToken)
}

func (h *accountHandler) CreateAccount(c echo.Context) error {
	var req dto.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(echo.ErrBadRequest.Code, dto.ErrorData{
			Code:    400,
			Message: "Invalid request body",
		})
	}

	if err := req.Validate(); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(echo.ErrBadRequest.Code, dto.ValidationErrorResponse{
				Code:    400,
				Message: "Validation failed",
				Errors:  validationErrors,
			})
		}
		return c.JSON(echo.ErrBadRequest.Code, dto.ErrorData{
			Code:    400,
			Message: err.Error(),
		})
	}

	if err := h.accountService.CreateAccount(c.Request().Context(), req); err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, dto.ErrorData{
			Code:    500,
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusOK)
}

func (h *accountHandler) AuthenticateAccount(c echo.Context) error {
	var req dto.AuthenticateAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(echo.ErrBadRequest.Code, dto.ErrorData{
			Code:    400,
			Message: "Invalid request body",
		})
	}

	if err := req.Validate(); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return c.JSON(echo.ErrBadRequest.Code, dto.ValidationErrorResponse{
				Code:    400,
				Message: "Validation failed",
				Errors:  validationErrors,
			})
		}
		return c.JSON(echo.ErrBadRequest.Code, dto.ErrorData{
			Code:    400,
			Message: err.Error(),
		})
	}

	token, err := h.accountService.AuthenticateAccount(c.Request().Context(), req)
	if err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, dto.ErrorData{
			Code:    500,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.AuthenticateAccountResponse{
		Token: token,
	})
}

func (h *accountHandler) GetAccountByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(echo.ErrBadRequest.Code, dto.ErrorData{
			Code:    400,
			Message: "Invalid account ID",
		})
	}

	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		return c.JSON(echo.ErrBadRequest.Code, dto.ErrorData{
			Code:    400,
			Message: "Invalid account ID: must be a positive number",
		})
	}

	account, err := h.accountService.GetAccountByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, dto.ErrorData{
			Code:    500,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.StandardResponse{
		Data: account,
	})
}

func (h *accountHandler) GetAccountByToken(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(echo.ErrBadRequest.Code, dto.ErrorData{
			Code:    400,
			Message: "Missing authorization header",
		})
	}

	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	account, err := h.accountService.GetAccountByToken(c.Request().Context(), token)
	if err != nil {
		return c.JSON(echo.ErrInternalServerError.Code, dto.ErrorData{
			Code:    500,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.StandardResponse{
		Data: account,
	})
}
