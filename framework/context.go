package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Context struct {
	request        *http.Request
	responseWriter http.ResponseWriter
	context        context.Context
	handler        ControllerHandler

	hasTimeout  bool
	writerMutex *sync.Mutex
}

func NewContext(r *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		request:        r,
		responseWriter: w,
		context:        r.Context(),
		writerMutex:    &sync.Mutex{},
	}
}

func (context *Context) WriterMutex() *sync.Mutex {
	return context.writerMutex
}

func (context *Context) GetRequest() *http.Request {
	return context.request
}

func (context *Context) GetResponse() http.ResponseWriter {
	return context.responseWriter
}

func (context *Context) SetHandler(handler ControllerHandler) {
	context.handler = handler
}

func (context *Context) SetHasTimeout() {
	context.hasTimeout = true
}

func (context *Context) HasTimeout() bool {
	return context.hasTimeout
}

func (context *Context) BaseContext() context.Context {
	return context.request.Context()
}

func (context *Context) Deadline() (deadline time.Time, ok bool) {
	return context.BaseContext().Deadline()
}

func (context *Context) Done() <-chan struct{} {
	return context.BaseContext().Done()
}

func (context *Context) Error() error {
	return context.BaseContext().Err()
}

func (context *Context) Value(key interface{}) interface{} {
	return context.BaseContext().Value(key)
}

// Feature: Parse the query of url !

func (context *Context) QueryAll() map[string][]string {
	if context.request != nil {
		return map[string][]string(context.request.URL.Query())
	}

	return map[string][]string{}
}

func (context *Context) QueryInt(key string, def int) int {
	params := context.QueryAll()

	if values, ok := params[key]; ok {
		len := len(values)
		if len > 0 {
			i, err := strconv.Atoi(values[len-1])
			if err != nil {
				return def
			}

			return i
		}
	}

	return def
}

func (context *Context) QueryString(key string, def string) string {
	params := context.QueryAll()

	if values, ok := params[key]; ok {
		len := len(values)
		if len > 0 {
			return values[len-1]
		}
	}

	return def
}

func (context *Context) QueryArray(key string, def []string) []string {
	params := context.QueryAll()

	if values, ok := params[key]; ok {
		return values
	}

	return def
}

// Feature: Parse the form of url !

func (context *Context) FormAll() map[string][]string {
	if context.request != nil {
		return map[string][]string(context.request.PostForm)
	}

	return map[string][]string{}
}

func (context *Context) FormInt(key string, def int) int {
	params := context.FormAll()

	if values, ok := params[key]; ok {
		len := len(values)
		if len > 0 {
			i, err := strconv.Atoi(values[len-1])
			if err != nil {
				return def
			}

			return i
		}
	}

	return def
}

func (context *Context) FormString(key string, def string) string {
	params := context.FormAll()

	if values, ok := params[key]; ok {
		len := len(values)
		if len > 0 {
			return values[len-1]
		}
	}

	return def
}

func (context *Context) FormArray(key string, def []string) []string {
	params := context.FormAll()

	if values, ok := params[key]; ok {
		return values
	}

	return def
}

// Feature: request json handle function

func (context *Context) BindJson(object interface{}) error {
	if context.request != nil {
		body, err := ioutil.ReadAll(context.request.Body)

		context.request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		err = json.Unmarshal(body, object)

		if err != nil {
			return err
		}
	} else {
		return errors.New("context.request empty !")
	}

	return nil
}

// Feature: response function

func (context *Context) Json(status int, object interface{}) error {
	if context.HasTimeout() {
		return nil
	}

	context.responseWriter.Header().Set("Content-Type", "application/json")
	context.responseWriter.WriteHeader(status)
	bytes, err := json.Marshal(object)
	if err != nil {
		context.responseWriter.WriteHeader(500)
		return err
	}

	context.responseWriter.Write(bytes)

	return nil
}

func (context *Context) HTML(status int, object interface{}, template string) error {
	return nil
}

func (context *Context) Text(status int, object interface{}) error {
	return nil
}
