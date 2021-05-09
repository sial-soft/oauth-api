package rest

import (
	"github.com/federicoleon/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"example@gmail.com","password":"the-password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid rest client when trying to login", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"example@gmail.com","password":"the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"Message": "invalid login credentials", "status": "404", "error": "not_found"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when try", err.Message)
}

func TestLoginUserInvalidCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"example@gmail.com","password":"the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"Message": "invalid login credentials", "status": 404, "error": "not_found"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJson(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"example@gmail.com","password":"the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody: `{
						"id": "1",
						"first_name": "Test",
						"last_name": "User",
						"email": "sample@gmail.com"
					}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshal users response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8080/users/login",
		ReqBody:      `{"email":"example@gmail.com","password":"the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody: `{
						"id": 1,
						"first_name": "Test",
						"last_name": "User",
						"email": "sample@gmail.com"
					}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 1, user.Id)
	assert.EqualValues(t, "Test", user.FirstName)
	assert.EqualValues(t, "User", user.LastName)
	assert.EqualValues(t, "sample@gmail.com", user.Email)
}
