package vsocket

import (
	"net"
	"reflect"
)

// DEFAULTACTION 默认
// BEFORACTION
// AFTERACTION
const (
	DEFAULTACTION = "Default"
	BEFORACTION   = "BeforeRequest"
	AFTERACTION   = "AfterRequest"
)

type module interface {
	Default(fd uint32, data map[string]string) bool
	BeforeRequest(fd uint32, data map[string]string) bool
	AfterRequest(fd uint32, data map[string]string) bool
}

type events interface {
	OnHandel(fd uint32, conn net.Conn) bool
	OnClose(fd uint32)
	OnMessage(fd uint32, msg map[string]string) bool
}

type routerMap struct {
	pools    map[string]func(uint32, map[string]string) bool
	strPools map[string]map[string]func(uint32, map[string]string) bool
	structs  map[string]module
	events   events
}

func newRouterMap() *routerMap {
	return &routerMap{
		pools:    make(map[string]func(uint32, map[string]string) bool),
		strPools: make(map[string]map[string]func(uint32, map[string]string) bool),
		structs:  make(map[string]module),
	}
}

//注册事件
func (r *routerMap) RegisterEvent(events events) {
	r.events = events
}

//注册单个逻辑
func (r *routerMap) RegisterFun(methodName string, funcs func(uint32, map[string]string) bool) bool {
	if _, exit := r.pools[methodName]; !exit {
		r.pools[methodName] = funcs
		return true
	}
	return false
}

// 结构体 注册
func (r *routerMap) RegisterStructFun(moduleName string, mod module) bool {
	if _, exit := r.strPools[moduleName]; exit {
		return false
	}
	r.strPools[moduleName] = make(map[string]func(uint32, map[string]string) bool)
	r.structs[moduleName] = mod

	temType := reflect.TypeOf(mod)
	temValue := reflect.ValueOf(mod)
	for i := 0; i < temType.NumMethod(); i++ {
		tem := temValue.Method(i).Interface()
		if temFunc, ok := tem.(func(uint32, map[string]string) bool); ok {
			r.strPools[moduleName][temType.Method(i).Name] = temFunc
		}
	}
	return true
}

func (r *routerMap) HookAction(funcionName string, fd uint32, data map[string]string) bool {
	if action, exit := r.pools[funcionName]; exit {
		return action(fd, data)
	} else {
		return false
	}
}

func (r *routerMap) HookModule(mouleName string, method string, fd uint32, data map[string]string) bool {
	if _, exit := r.strPools[mouleName]; !exit {
		return false
	}

	if r.strPools[mouleName][BEFORACTION](fd, data) == false {
		return false
	}
	if action, exit := r.strPools[mouleName][method]; exit {
		if action(fd, data) == false {
			return false
		}
	} else {
		if r.strPools[mouleName][DEFAULTACTION](fd, data) == false {
			return false
		}
	}
	if r.strPools[mouleName][AFTERACTION](fd, data) == false {
		return false
	}
	return true
}

func (r *routerMap) OnClose(fd uint32) {
	if r.events != nil {
		r.events.OnClose(fd)
	}
}

func (r *routerMap) OnHandel(fd uint32, conn net.Conn) bool {
	if r.events != nil {
		return r.events.OnHandel(fd, conn)
	}
	return true
}

func (r *routerMap) OnMessage(fd uint32, msg map[string]string) bool {
	if r.events != nil {
		return r.events.OnMessage(fd, msg)
	}
	return true
}
