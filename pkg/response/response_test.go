package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AccelByte/go-restful-plugins/v4/pkg/logger/event"
	"github.com/AccelByte/go-restful-plugins/v4/pkg/util"
	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
)

// nolint: dupl // most part of the test is identical
func TestWriteSuccess(t *testing.T) {
	t.Parallel()

	ws := new(restful.WebService)
	ws.Filter(event.Log("test", "go-restful-plugins", util.ExtractDefault))

	type ResponseTest struct {
		Message string `json:"message"`
	}

	ws.Route(
		ws.GET("/namespace/{namespace}/user/{id}").
			Param(restful.PathParameter("namespace", "namespace")).
			Param(restful.PathParameter("id", "user ID")).
			To(func(request *restful.Request, response *restful.Response) {
				Write(request, response, http.StatusOK, 0, 30, "test success", ResponseTest{Message: "success"})
			}))

	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest(http.MethodGet, "/namespace/abc/user/def", nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")

	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	decoder := json.NewDecoder(resp.Body)

	var responseTest ResponseTest
	_ = decoder.Decode(&responseTest)

	assert.Equal(t, http.StatusOK, resp.Code, "response status code should be %v", http.StatusOK)
	assert.Equal(t, ResponseTest{Message: "success"}, responseTest,
		"response body must be %+v", ResponseTest{Message: "success"})
}

// nolint: dupl // most part of the test is identical
func TestWriteErrorWarning(t *testing.T) {
	t.Parallel()

	ws := new(restful.WebService)
	ws.Filter(event.Log("test", "go-restful-plugins", util.ExtractDefault))

	ws.Route(
		ws.GET("/namespace/{namespace}/user/{id}").
			Param(restful.PathParameter("namespace", "namespace")).
			Param(restful.PathParameter("id", "user ID")).
			To(func(request *restful.Request, response *restful.Response) {
				WriteError(request, response, http.StatusNotFound, 0, errors.New("123"), &Error{
					ErrorCode:    30,
					ErrorMessage: "abc",
					ErrorLogMsg:  "cba",
				})
			}))

	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest(http.MethodGet, "/namespace/abc/user/def", nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")

	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	decoder := json.NewDecoder(resp.Body)

	var responseTest Error
	_ = decoder.Decode(&responseTest)

	expected := Error{
		ErrorCode:    30,
		ErrorMessage: "abc",
		ErrorLogMsg:  "",
	}

	assert.Equal(t, http.StatusNotFound, resp.Code, "response status code should be %v", http.StatusOK)
	assert.Equal(t, expected, responseTest, "response body must be %+v", expected)
}

// nolint: dupl // most part of the test is identical
func TestWriteErrorInternalServerError(t *testing.T) {
	t.Parallel()

	ws := new(restful.WebService)
	ws.Filter(event.Log("test", "go-restful-plugins", util.ExtractDefault))

	ws.Route(
		ws.GET("/namespace/{namespace}/user/{id}").
			Param(restful.PathParameter("namespace", "namespace")).
			Param(restful.PathParameter("id", "user ID")).
			To(func(request *restful.Request, response *restful.Response) {
				WriteError(request, response, http.StatusInternalServerError, 0, errors.New("323"), &Error{
					ErrorCode:    31,
					ErrorMessage: "abc",
					ErrorLogMsg:  "cba",
				})
			}))

	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest(http.MethodGet, "/namespace/abc/user/def", nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")

	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	decoder := json.NewDecoder(resp.Body)

	var responseTest Error
	_ = decoder.Decode(&responseTest)

	expected := Error{
		ErrorCode:    31,
		ErrorMessage: "abc",
		ErrorLogMsg:  "",
	}

	assert.Equal(t, http.StatusInternalServerError, resp.Code, "response status code should be %v", http.StatusOK)
	assert.Equal(t, expected, responseTest, "response body must be %+v", expected)
}

// nolint: dupl // most part of the test is identical
func TestWriteErrorWithEventIDWarning(t *testing.T) {
	t.Parallel()

	ws := new(restful.WebService)
	ws.Filter(event.Log("test", "go-restful-plugins", util.ExtractDefault))

	ws.Route(
		ws.GET("/namespace/{namespace}/user/{id}").
			Param(restful.PathParameter("namespace", "namespace")).
			Param(restful.PathParameter("id", "user ID")).
			To(func(request *restful.Request, response *restful.Response) {
				WriteErrorWithEventID(request, response, http.StatusNotFound, 0, 40, errors.New("123"),
					&Error{
						ErrorCode:    30,
						ErrorMessage: "abc",
						ErrorLogMsg:  "cba",
					})
			}))

	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest(http.MethodGet, "/namespace/abc/user/def", nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")

	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	decoder := json.NewDecoder(resp.Body)

	var responseTest Error
	_ = decoder.Decode(&responseTest)

	expected := Error{
		ErrorCode:    30,
		ErrorMessage: "abc",
		ErrorLogMsg:  "",
	}

	assert.Equal(t, http.StatusNotFound, resp.Code, "response status code should be %v", http.StatusOK)
	assert.Equal(t, expected, responseTest, "response body must be %+v", expected)
}

func TestWriteErrorWithEventIDInternalServerError(t *testing.T) {
	t.Parallel()

	ws := new(restful.WebService)
	ws.Filter(event.Log("test", "go-restful-plugins", util.ExtractDefault))

	ws.Route(
		ws.GET("/namespace/{namespace}/user/{id}").
			Param(restful.PathParameter("namespace", "namespace")).
			Param(restful.PathParameter("id", "user ID")).
			To(func(request *restful.Request, response *restful.Response) {
				WriteErrorWithEventID(request, response, http.StatusInternalServerError, 0, 41, errors.New("323"),
					&Error{
						ErrorCode:    31,
						ErrorMessage: "abc",
						ErrorLogMsg:  "cba",
					})
			}))

	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest(http.MethodGet, "/namespace/abc/user/def", nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")

	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	decoder := json.NewDecoder(resp.Body)

	var responseTest Error
	_ = decoder.Decode(&responseTest)

	expected := Error{
		ErrorCode:    31,
		ErrorMessage: "abc",
		ErrorLogMsg:  "",
	}

	assert.Equal(t, http.StatusInternalServerError, resp.Code, "response status code should be %v", http.StatusOK)
	assert.Equal(t, expected, responseTest, "response body must be %+v", expected)
}
