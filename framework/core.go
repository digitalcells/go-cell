package framework

import (
	"log"
	"net/http"
	"strings"
)

type Core struct {
	router map[string]*Tree
}

func NewCore() *Core {
	router := map[string]*Tree{}
	router["GET"] = NewTree()
	router["POST"] = NewTree()
	router["PUT"] = NewTree()
	router["DELETE"] = NewTree()
	return &Core{router: router}
}

func (core *Core) Get(url string, handler ControllerHandler) {
	if err := core.router["GET"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (core *Core) Post(url string, handler ControllerHandler) {
	if err := core.router["POST"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (core *Core) Put(url string, handler ControllerHandler) {
	if err := core.router["PUT"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (core *Core) Delete(url string, handler ControllerHandler) {
	if err := core.router["DELETE"].AddRouter(url, handler); err != nil {
		log.Fatal("add router error: ", err)
	}
}

func (core *Core) Group(prefix string) IGroup {
	return NewGroup(core, prefix)
}

func (core *Core) FindRouteByRequest(request *http.Request) ControllerHandler {
	uri := request.URL.Path
	method := request.Method

	_method := strings.ToUpper(method)

	if handlers, ok := core.router[_method]; ok {
		return handlers.FindHandler(uri)
	}

	return nil
}

func (core *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Println("Core:Http:Server")

	context := NewContext(request, response)

	router := core.FindRouteByRequest(request)

	if router == nil {
		return
	}

	log.Println("Core:Router")

	context.SetHandler(router)

	router(context)
}
