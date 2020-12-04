package ports

import (
	"errors"
	"net/http"

	"github.com/crounch-me/back/internal"
	"github.com/crounch-me/back/internal/account/app"
	"github.com/crounch-me/back/internal/common/server"
	"github.com/crounch-me/back/util"
	"github.com/gin-gonic/gin"
)

const (
	loginPath  = "/account/login"
	signupPath = "/account/signup"
	userPath   = "/users"
	mePath     = "/me"
	logoutPath = "/logout"
)

type GinServer struct {
	accountService *app.AccountService
	validator      *util.Validator
}

func NewGinServer(accountService *app.AccountService, validator *util.Validator) (*GinServer, error) {
	if accountService == nil {
		return nil, errors.New("account gin server accountService is nil")
	}

	if validator == nil {
		return nil, errors.New("account gin server validator is nil")
	}

	return &GinServer{
		accountService: accountService,
		validator:      validator,
	}, nil
}

func (s *GinServer) ConfigureRoutes(r *gin.Engine) {
	r.POST(signupPath, s.Signup)
	r.OPTIONS(signupPath, server.OptionsHandler([]string{http.MethodPost}))

	r.POST(loginPath, s.Login)
	r.OPTIONS(loginPath, server.OptionsHandler([]string{http.MethodPost}))

	r.POST(logoutPath, s.Logout)
	r.OPTIONS(logoutPath, server.OptionsHandler([]string{http.MethodPost}))
}

// Signup creates a new user with his email and password
// @Summary Creates a new user with his email and password
// @ID signup
// @Tags account
// @Accept json
// @Param user body SignupRequest true "User to signup with"
// @Success 201
// @Failure 400 {object} internal.Error
// @Failure 500 {object} internal.Error
// @Router /account/signup [post]
func (s *GinServer) Signup(c *gin.Context) {
	signupRequest := &SignupRequest{}

	err := server.UnmarshalPayload(c.Request.Body, signupRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.UnmarshalErrorCode))
		return
	}

	err = s.validator.Struct(signupRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.InvalidErrorCode))
		return
	}

	err = s.accountService.Signup(signupRequest.Email, signupRequest.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, internal.NewError(internal.UnknownErrorCode))
		return
	}

	c.Status(http.StatusCreated)
}

// Login creates a new user authorization when email is found and password is valid
// @Summary Creates a new user authorization when email is found and password is valid
// @ID login
// @Tags account
// @Accept json
// @Produce  json
// @Param user body LoginRequest true "User to login with"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} internal.Error
// @Failure 500 {object} internal.Error
// @Router /account/login [post]
func (s *GinServer) Login(c *gin.Context) {
	loginRequest := &LoginRequest{}

	err := server.UnmarshalPayload(c.Request.Body, loginRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.UnmarshalErrorCode))
		return
	}

	err = s.validator.Struct(loginRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.InvalidErrorCode))
		return
	}

	token, err := s.accountService.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, internal.NewError(internal.UnknownErrorCode))
		return
	}

	tokenResponse := &TokenResponse{
		Token: token,
	}

	server.JSON(c, tokenResponse)
}

// Logout removes the user authorization when the user token is found
// @Summary Removes the user authorization when the user token is found
// @ID logout
// @Tags account
// @Success 204
// @Failure 403 {object} internal.Error
// @Failure 500 {object} internal.Error
// @Security ApiKeyAuth
// @Router /account/logout [post]
func (s *GinServer) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if token == "" {
		c.Status(http.StatusNoContent)
		return
	}

	userUUID, err := s.accountService.GetUserUUIDByToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, internal.NewError(internal.ForbiddenErrorCode))
		return
	}

	err = s.accountService.Logout(userUUID, token)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}