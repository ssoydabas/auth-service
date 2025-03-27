package handler

import (
	"net/http"
	"strconv"

	"github.com/ssoydabas/auth-service/internal/dto"
	"github.com/ssoydabas/auth-service/internal/service"
	"github.com/ssoydabas/auth-service/pkg/errors"
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
	SetResetPasswordToken(c echo.Context) error
	ResetPassword(c echo.Context) error
	GetAccountEmailVerificationTokenByID(c echo.Context) error
	VerifyAccountEmail(c echo.Context) error
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
	e.POST("/accounts/set-reset-password-token", h.SetResetPasswordToken)
	e.POST("/accounts/reset-password", h.ResetPassword)
	e.GET("/accounts/get-email-verification-token/:id", h.GetAccountEmailVerificationTokenByID)
	e.POST("/accounts/verify-email", h.VerifyAccountEmail)
}

// @Summary Create a new account
// @Description Create a new account
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body dto.CreateAccountRequest true "Account details"
// @Success 201 {object} nil "Account created successfully"
// @Failure 400 {object} dto.ValidationErrorResponse
// @Failure 409 {object} dto.ErrorData
// @Failure 500 {object} dto.ErrorData
// @Router /accounts [post]
func (h *accountHandler) CreateAccount(c echo.Context) error {
	var req dto.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return errors.BadRequestError("Invalid request body")
	}

	if err := req.Validate(); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.ValidationError("Validation failed", validationErrors)
		}
		return errors.BadRequestError(err.Error())
	}

	token, err := h.accountService.CreateAccount(c.Request().Context(), req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.InternalError(err)
	}

	return c.JSON(http.StatusCreated, dto.VerificationCodeResponse{
		AuthenticateAccountResponse: dto.AuthenticateAccountResponse{
			Token: token,
		},
		VerificationCode: token,
	})
}

// @Summary Authenticate an account
// @Description Authenticate an account and receive a JWT token
// @Tags accounts
// @Accept json
// @Produce json
// @Param account body dto.AuthenticateAccountRequest true "Account credentials"
// @Success 200 {object} dto.AuthenticateAccountResponse
// @Failure 400 {object} dto.ValidationErrorResponse
// @Failure 401 {object} dto.ErrorData
// @Failure 500 {object} dto.ErrorData
// @Router /accounts/authenticate [post]
func (h *accountHandler) AuthenticateAccount(c echo.Context) error {
	var req dto.AuthenticateAccountRequest
	if err := c.Bind(&req); err != nil {
		return errors.BadRequestError("Invalid request body")
	}

	if err := req.Validate(); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.ValidationError("Validation failed", validationErrors)
		}
		return errors.BadRequestError(err.Error())
	}

	token, err := h.accountService.AuthenticateAccount(c.Request().Context(), req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.InternalError(err)
	}

	return c.JSON(http.StatusOK, dto.AuthenticateAccountResponse{
		Token: token,
	})
}

// @Summary Get an account by ID
// @Description Get an account details by their ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path integer true "Account ID"
// @Success 200 {object} dto.StandardResponse{data=dto.AccountResponse}
// @Failure 400 {object} dto.ErrorData
// @Failure 404 {object} dto.ErrorData
// @Failure 500 {object} dto.ErrorData
// @Router /accounts/{id} [get]
func (h *accountHandler) GetAccountByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return errors.BadRequestError("Account ID is required")
	}

	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		return errors.BadRequestError("Invalid account ID: must be a positive number")
	}

	account, err := h.accountService.GetAccountByID(c.Request().Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.InternalError(err)
	}

	return c.JSON(http.StatusOK, dto.StandardResponse{
		Data: account,
	})
}

// @Summary Get current account details
// @Description Get account details using JWT token
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.AccountResponse
// @Failure 400 {object} dto.ErrorData
// @Failure 401 {object} dto.ErrorData
// @Failure 500 {object} dto.ErrorData
// @Router /accounts/me [get]
func (h *accountHandler) GetAccountByToken(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return errors.BadRequestError("Missing authorization header")
	}

	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	account, err := h.accountService.GetAccountByToken(c.Request().Context(), token)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.InternalError(err)
	}

	return c.JSON(http.StatusOK, mapAccountToResponse(*account))
}

// @Summary Request password reset
// @Description Request a password reset token to be sent to email, either email or phone number is required
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body dto.SetResetPasswordTokenRequest true "Email address"
// @Success 200 {object} dto.StandardResponse{data=string} "Reset token"
// @Failure 400 {object} dto.ValidationErrorResponse
// @Failure 404 {object} dto.ErrorData
// @Failure 500 {object} dto.ErrorData
// @Router /accounts/set-reset-password-token [post]
func (h *accountHandler) SetResetPasswordToken(c echo.Context) error {
	var req dto.SetResetPasswordTokenRequest
	if err := c.Bind(&req); err != nil {
		return errors.BadRequestError("Invalid request body")
	}

	if err := req.Validate(); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.ValidationError("Validation failed", validationErrors)
		}
		return errors.BadRequestError(err.Error())
	}

	token, err := h.accountService.SetResetPasswordToken(c.Request().Context(), req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.InternalError(err)
	}

	return c.JSON(http.StatusOK, dto.StandardResponse{
		Data: token,
	})
}

// @Summary Reset password
// @Description Reset password using the token received via email
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} nil "Password reset successfully"
// @Failure 400 {object} dto.ValidationErrorResponse
// @Failure 401 {object} dto.ErrorData
// @Failure 500 {object} dto.ErrorData
// @Router /accounts/reset-password [post]
func (h *accountHandler) ResetPassword(c echo.Context) error {
	var req dto.ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return errors.BadRequestError("Invalid request body")
	}

	if err := req.Validate(); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.ValidationError("Validation failed", validationErrors)
		}
		return errors.BadRequestError(err.Error())
	}

	if err := h.accountService.ResetPassword(c.Request().Context(), req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.InternalError(err)
	}

	return c.NoContent(http.StatusOK)
}

// @Summary Get email verification token by account ID
// @Description Get email verification token by account ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path integer true "Account ID"
// @Success 200 {object} dto.StandardResponse{data=string} "Verification token"
// @Failure 400 {object} dto.ErrorData
// @Failure 404 {object} dto.ErrorData
// @Failure 500 {object} dto.ErrorData
// @Router /accounts/email-verification-token/{id} [get]
func (h *accountHandler) GetAccountEmailVerificationTokenByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return errors.BadRequestError("Account ID is required")
	}

	if _, err := strconv.ParseUint(id, 10, 64); err != nil {
		return errors.BadRequestError("Invalid account ID: must be a positive number")
	}

	token, err := h.accountService.GetAccountEmailVerificationTokenByID(c.Request().Context(), id)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.InternalError(err)
	}

	return c.JSON(http.StatusOK, dto.StandardResponse{
		Data: token,
	})
}

// @Summary Verify email address
// @Description Verify account email using verification token
// @Tags accounts
// @Accept json
// @Produce json
// @Param request body dto.VerifyAccountRequest true "Verification token"
// @Success 200 {object} nil "Email verified successfully"
// @Failure 400 {object} dto.ValidationErrorResponse
// @Failure 401 {object} dto.ErrorData
// @Failure 500 {object} dto.ErrorData
// @Router /accounts/verify-email [post]
func (h *accountHandler) VerifyAccountEmail(c echo.Context) error {
	var req dto.VerifyAccountRequest
	if err := c.Bind(&req); err != nil {
		return errors.BadRequestError("Invalid request body")
	}

	if err := req.Validate(); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return errors.ValidationError("Validation failed", validationErrors)
		}
		return errors.BadRequestError(err.Error())
	}

	if err := h.accountService.VerifyAccountEmail(c.Request().Context(), req); err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return appErr
		}
		return errors.InternalError(err)
	}

	return c.NoContent(http.StatusOK)
}
