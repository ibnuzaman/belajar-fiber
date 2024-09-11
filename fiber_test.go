package main

import (
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var app = fiber.New()

func TestFiber(t *testing.T) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response, err := app.Test(request)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200")

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "Hello, World!", string(bytes), "Body should be \"Hello, World!\"")
}

func TestCtx(t *testing.T) {

	app.Get("/hello", func(c *fiber.Ctx) error {
		name := c.Query("name", "Guest")
		return c.SendString("Hello " + name)
	})

	request := httptest.NewRequest(http.MethodGet, "/hello?name=ibnu", nil)
	response, err := app.Test(request)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200")

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "Hello ibnu", string(bytes), "Body should be \"Hello, World!\"")

	request = httptest.NewRequest(http.MethodGet, "/hello?ibnu", nil)
	response, err = app.Test(request)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200")

	bytes, err = io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "Hello Guest", string(bytes), "Body should be \"Hello, World!\"")
}

func TestHttpRequest(t *testing.T) {
	app.Get("/req", func(c *fiber.Ctx) error {
		first := c.Get("firstname")
		last := c.Cookies("lastname")
		return c.SendString("hello " + first + " " + last)
	})

	request := httptest.NewRequest(http.MethodGet, "/req", nil)
	request.Header.Set("firstname", "Ibnu")
	request.AddCookie(&http.Cookie{Name: "lastname", Value: "zaman"})
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello Ibnu zaman", string(bytes), "Body should be \"fisrtname\"")
}

func TestParameter(t *testing.T) {
	app.Get("/user/:userId/order/:orderId", func(c *fiber.Ctx) error {
		userId := c.Params("userId")
		orderId := c.Params("orderId")
		return c.SendString(userId + " " + orderId)
	})

	request := httptest.NewRequest(http.MethodGet, "/user/123/order/456", nil)
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200")
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "123 456", string(bytes), "Body should be \"123 456\"")

}

func TestNewRequestForm(t *testing.T) {
	app.Post("/post", func(c *fiber.Ctx) error {
		name := c.FormValue("name")
		return c.SendString("Hello " + name)
	})

	body := strings.NewReader("name=Ibnu")
	request := httptest.NewRequest("POST", "/post", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "Hello Ibnu", string(bytes))
}

var contohFile []byte

func TestMultiPartForm(t *testing.T) {
	app.Post("/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		err = c.SaveFile(file, file.Filename)
		if err != nil {
			return err
		}
		return c.SendString("upload sukses")
	})

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, _ := writer.CreateFormFile("file", "contoh.txt")
	file.Write(contohFile)
	writer.Close()

	request := httptest.NewRequest(http.MethodPost, "/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200")
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, "upload sukses", string(bytes))
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestLogin(t *testing.T) {
	app.Post("/login", func(c *fiber.Ctx) error {
		body := c.Body()
		request := new(LoginRequest)
		err := json.Unmarshal(body, &request)
		if err != nil {
			return err
		}
		return c.SendString("Login sukses" + request.Username)
	})

	body := strings.NewReader(`{"username":"ibnu", "password":"123456"}`)
	request := httptest.NewRequest(http.MethodPost, "/login", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200", string(bytes))
}

type RegisterRequest struct {
	Username string `json:"username" xml:"username" form:"username"`
	Password string `json:"password" xml:"password" form:"password"`
	Name     string `json:"name" xml:"name" form:"name"`
}

func TestBodyParser(t *testing.T) {
	app.Post("/register", func(c *fiber.Ctx) error {
		body := new(RegisterRequest)
		err := c.BodyParser(body)
		if err != nil {
			return err
		}
		return c.SendString("register sukses" + body.Username)
	})
}

func TestBodyParserJSON(t *testing.T) {
	TestBodyParser(t)
	body := strings.NewReader(`{"username":"ibnu", "password":"123456", "name" : "ibnu"}`)
	request := httptest.NewRequest(http.MethodPost, "/register", body)
	request.Header.Set("Content-Type", "application/json")
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200", string(bytes))
}

func TestBodyParserForm(t *testing.T) {
	TestBodyParser(t)
	body := strings.NewReader(`username=Ibnu&password=123456&name=Ibnu`)
	request := httptest.NewRequest(http.MethodPost, "/register", body)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200", string(bytes))
}

func TestBodyParserXML(t *testing.T) {
	TestBodyParser(t)
	body := strings.NewReader(
		`<RegisterRequest> 
  			  <Username>Ibnu</Username>
			  <password>123456</password>
			  <name>Ibnu</name>
			</RegisterRequest>`)
	request := httptest.NewRequest(http.MethodPost, "/register", body)
	request.Header.Set("Content-Type", "application/xml")
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200", string(bytes))
}

func TestResponseJSON(t *testing.T) {
	app.Get("/user", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"username": "Ibnu",
			"password": "123456",
			"name":     "Ibnu",
		})
	})

	request := httptest.NewRequest(http.MethodGet, "/user", nil)
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, `{"name":"Ibnu","password":"123456","username":"Ibnu"}`, string(bytes))

}

func TestDownloadFile(t *testing.T) {
	app.Get("/download", func(c *fiber.Ctx) error {
		return c.Download("contoh.txt", "contoh.txt")
	})

	request := httptest.NewRequest(http.MethodGet, "/download", nil)
	response, err := app.Test(request)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200")
	assert.Equal(t, "attachment; filename=\"contoh.txt\"", response.Header.Get("Content-Disposition"))
	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, fiber.StatusOK, response.StatusCode, "Status should be 200", string(bytes))
}

func TestClient(t *testing.T) {
	client := fiber.AcquireClient()

	agent := client.Get("https://example.com")
	status, response, errors := agent.String()
	assert.Nil(t, errors)
	assert.Equal(t, status, fiber.StatusOK)
	//assert.Contains(t, response, "Example Domain")
	assert.Contains(t, strings.ToLower(response), "example domain")
	fiber.ReleaseClient(client)
}
