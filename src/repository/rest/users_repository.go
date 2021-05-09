package rest

import (
	"encoding/json"
	"github.com/federicoleon/golang-restclient/rest"
	"github.com/sial-soft/oauth-api/src/domain/users"
	"github.com/sial-soft/oauth-api/src/utils/rest_errors"
	"time"
)

var (
	restClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080",
		Timeout: 100 * time.Millisecond,
	}
)

type RestUsersRepository interface {
	LoginUser(string, string) (*users.User, *rest_errors.RestErr)
}

type usersRepository struct {
}

func NewRestUserRepository() RestUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) LoginUser(email string, password string) (*users.User, *rest_errors.RestErr) {
	request := users.UserLoginRequest{
		Email:    email,
		Password: password,
	}
	response := restClient.Post("/users/login", request)
	if response == nil || response.Response == nil {
		return nil, rest_errors.NewInternalError("invalid rest client when trying to login")
	}
	if response.StatusCode > 299 {
		var restErr rest_errors.RestErr
		err := json.Unmarshal(response.Bytes(), &restErr)
		if err != nil {
			return nil, rest_errors.NewInternalError("invalid error interface when try")
		}
		return nil, &restErr
	}
	var user users.User
	if err := json.Unmarshal(response.Bytes(), &user); err != nil {
		return nil, rest_errors.NewInternalError("error when trying to unmarshal users response")
	}
	return &user, nil
}
