package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/crounch-me/back/configuration"
	"github.com/crounch-me/back/domain"
	"github.com/crounch-me/back/domain/authorization.go"
	"github.com/crounch-me/back/domain/lists"
	"github.com/crounch-me/back/domain/products"
	"github.com/crounch-me/back/domain/users"
	"github.com/crounch-me/back/storage"
	"github.com/crounch-me/back/storage/mock"
	"github.com/crounch-me/back/storage/postgres"
)

const (
	// ContextUserID is the key to retrieve user id from gin context
	ContextUserID = "ContextUserID"
)

type Services struct {
	Authorization *authorization.AuthorizationService
	List          *lists.ListService
	Product       *products.ProductService
	User          *users.UserService
}

// Context holds everything to respond to requests
type Context struct {
	Storage  storage.Storage
	Config   *configuration.Config
	Validate *validator.Validate
	Services *Services
}

// NewContext creates and initialize everything for the requests
func NewContext(config *configuration.Config) *Context {
	var storage storage.Storage

	if config.Mock {
		storage = mock.NewStorageMock()
	} else {
		postgres.InitDB(config.DBConnectionURI)
		storage = postgres.NewStorage(config.DBConnectionURI, config.DBSchema)
	}

	return &Context{
		Storage:  storage,
		Config:   config,
		Validate: validator.New(),
		Services: NewServices(storage),
	}
}

// LogAndSendError logs and sends the error
func (hc *Context) LogAndSendError(c *gin.Context, err error) {
	var status int
	if err, ok := err.(*domain.Error); ok {
		status = hc.ErrorCodeToHTTPStatus(err.Code)
	} else {
		status = http.StatusInternalServerError
	}

	logrus.WithError(err).Error("an error occured")

	c.AbortWithStatusJSON(status, err)
}

// UnmarshalPayload unmarshal request payload from context into object
func (hc *Context) UnmarshalPayload(c *gin.Context, i interface{}) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(body), i)
	if err != nil {
		return err
	}
	return nil
}

func (hc *Context) GetUserIDFromContext(c *gin.Context) (string, *domain.Error) {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return "", domain.NewError(domain.UnknownErrorCode)
	}
	return userID.(string), nil
}

func (hc *Context) UnmarshalAndValidate(c *gin.Context, i interface{}) *domain.Error {
	err := hc.UnmarshalPayload(c, i)

	if err != nil {
		return domain.NewErrorWithCause(domain.UnmarshalErrorCode, err)
	}

	err = hc.Validate.Struct(i)
	if err != nil {
		return domain.NewErrorWithCause(domain.InvalidErrorCode, err)
	}

	return nil
}

// NewServices create an object which contains all necessary services
func NewServices(storage storage.Storage) *Services {
	return &Services{
		Authorization: &authorization.AuthorizationService{
			AuthorizationStorage: storage,
			UserStorage:          storage,
		},
		List: &lists.ListService{
			ListStorage:    storage,
			ProductStorage: storage,
		},
		Product: &products.ProductService{
			ProductStorage: storage,
		},
		User: &users.UserService{
			UserStorage: storage,
		},
	}
}
