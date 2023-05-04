package rest

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password":"the-password"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restClient response when trying to login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password":"the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials", "status": "404", "error": "not_found"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password":"the-password"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message":"invalid login credentials", "status": 404, "error": "not_found"}`,
	})
	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid login credentials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password":"the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "5","first_name": "facundo","last_name": "centeno","email": "facundo@gmail.com"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying to unmarshal de users login response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "https://api.bookstore.com/users/login",
		ReqBody:      `{"email":"email@gmail.com", "password":"the-password"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 5,"first_name": "facundo","last_name": "centeno","email": "facundo@gmail.com"}`,
	})

	repository := usersRepository{}
	user, err := repository.LoginUser("email@gmail.com", "the-password")
	fmt.Println(user)
	fmt.Println(err)

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 5, user.Id)
	assert.EqualValues(t, "facundo", user.FirstName)
	assert.EqualValues(t, "centeno", user.LastName)
	assert.EqualValues(t, "facundo@gmail.com", user.Email)
}
