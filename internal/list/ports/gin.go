package ports

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/crounch-me/back/internal"
	authorizationApp "github.com/crounch-me/back/internal/authorization/app"
	"github.com/crounch-me/back/internal/common/server"
	listApp "github.com/crounch-me/back/internal/list/app"
	userApp "github.com/crounch-me/back/internal/user/app"
	"github.com/crounch-me/back/util"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

const (
	listPath        = "/lists"
	listWithIDPath  = "/lists/:listID"
	archiveListPath = "/lists/:listID/archive"
)

type GinServer struct {
	authorizationService *authorizationApp.AuthorizationService
	listService          *listApp.ListService
	userService          *userApp.UserService
	validator            *util.Validator
}

func NewGinServer(listService *listApp.ListService, authorizationService *authorizationApp.AuthorizationService, validator *util.Validator) (*GinServer, error) {
	if listService == nil {
		return nil, errors.New("listService is nil")
	}

	if authorizationService == nil {
		return nil, errors.New("authorizationService is nil")
	}

	if validator == nil {
		return nil, errors.New("validator is nil")
	}

	return &GinServer{
		listService: listService,
		validator:   validator,
	}, nil
}

func (s *GinServer) ConfigureRoutes(r *gin.Engine) {
	r.POST(listPath, server.CheckUserAuthorization(s.authorizationService), s.CreateList)
	r.GET(listPath, server.CheckUserAuthorization(s.authorizationService), s.GetUserLists)
	r.OPTIONS(listPath, server.OptionsHandler([]string{http.MethodGet, http.MethodPost}))
}

func (s *GinServer) CreateList(c *gin.Context) {
	list := &CreateListRequest{}

	err := server.UnmarshalPayload(c.Request.Body, list)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.UnmarshalErrorCode))
		return
	}

	err = s.validator.Struct(list)
	if err != nil {
		fields := make([]*internal.FieldError, 0)
		for _, e := range err.(validator.ValidationErrors) {
			field := &internal.FieldError{
				Error: e.Tag(),
				Name:  e.Field(),
			}
			fields = append(fields, field)
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, internal.NewError(internal.InvalidErrorCode).WithFields(fields))
		return
	}

	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, internal.NewError(internal.ForbiddenErrorCode))
		return
	}

	listUUID, err := s.listService.CreateList(userUUID, list.Name)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, internal.NewError(internal.UnknownErrorCode))
		return
	}

	c.Header(server.HeaderContentLocation, "/lists/"+listUUID)
	c.Status(http.StatusCreated)
}

func (h *GinServer) GetUserLists(c *gin.Context) {
	userUUID, err := server.GetUserIDFromContext(c)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	lists, err := h.listService.GetUserLists(userUUID)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	listsResponse := make([]*List, 0)
	for _, list := range lists {
		products := make([]*Product, 0)
		for _, p := range list.Products() {
			product := &Product{
				UUID: p.UUID(),
			}
			products = append(products, product)
		}

		listResponse := &List{
			UUID:         list.UUID(),
			Name:         list.Name(),
			CreationDate: list.CreationDate(),
			Contributors: list.Contributors(),
			Products:     products,
		}

		listsResponse = append(listsResponse, listResponse)
	}

	server.JSON(c, listsResponse)
}
