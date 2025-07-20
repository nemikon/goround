package goround

import (
	"reflect"
	"sync"
)

type serviceConfig struct {
	enableStructInject   bool
	enableFunctionInject bool
	structTag            string
	structDefaultInject  bool
	functionInjectName   string
}

type service struct {
	instance    any
	preparingMu sync.RWMutex
	preparing   bool
	ready       bool
	readyMu     sync.RWMutex
	config      serviceConfig
}

type option func(config serviceConfig) serviceConfig

type Container struct {
	mu             sync.RWMutex
	services       map[string]*service
	defaultOptions []option
}

func NewContainer(defaultOptions ...option) *Container {
	return &Container{
		services:       make(map[string]*service),
		defaultOptions: defaultOptions,
	}
}

func RegisterBasicValue[T any](c *Container, value T) {
	Register(c, value, WithFunctionInject(false), WithStructInject(false))
}

func Register[T any](c *Container, empty T, options ...option) {
	c.mu.Lock()
	defer c.mu.Unlock()
	config := generateServiceConfig(c, options...)

	if config.enableFunctionInject {
		// verify function existis
		t := reflect.TypeOf(empty)
		_, ok := t.MethodByName(config.functionInjectName)
		if !ok {
			panic("empty doesn't have a function called " + config.functionInjectName)
		}
	}

	if config.enableFunctionInject || config.enableStructInject {
		// verify that empty is a pointer
		t := reflect.TypeOf(empty)
		if t.Kind() != reflect.Pointer {
			panic("empty must be a pointer when struct or function inject is enabled")
		}
	}

	name := getName[T]()
	service := &service{
		instance: empty,
		config:   config,
	}

	c.services[name] = service
}

func generateServiceConfig(c *Container, options ...option) serviceConfig {
	config := serviceConfig{
		structTag:          "goround",
		functionInjectName: "InjectDependencies",
	}
	for _, option := range c.defaultOptions {
		config = option(config)
	}

	for _, option := range options {
		config = option(config)
	}
	return config
}

func Get[I any](c *Container) I {
	name := getName[I]()
	s := GetByName(c, name)

	return s.(I)
}

func GetByName(c *Container, name string) any {
	service := getPreparingServiceByName(c, name)

	service.readyMu.RLock()
	defer service.readyMu.RUnlock()
	return service.instance
}

func getPreparingServiceByName(c *Container, name string) *service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, ok := c.services[name]
	if !ok {
		panic("nope")
	}

	prepareService(c, service)

	return service
}

func prepareService(c *Container, s *service) {
	s.preparingMu.Lock()
	if s.preparing {
		s.preparingMu.Unlock()
		return
	}

	s.preparing = true
	s.readyMu.Lock()
	defer s.readyMu.Unlock()
	s.preparingMu.Unlock()

	injectByFunc(c, s)
	injectByStruct(c, s)
	s.ready = true
}

func injectByStruct(c *Container, s *service) {
	if !s.config.enableStructInject {
		return
	}
	v := reflect.ValueOf(s.instance)
	st := v.Elem()
	stt := st.Type()

	for i := 0; i < stt.NumField(); i++ {
		field := stt.Field(i)

		if !field.IsExported() {
			continue
		}

		tag := field.Tag.Get(s.config.structTag)
		if tag == "ignore" {
			continue
		}
		if !s.config.structDefaultInject && tag != "inject" {
			continue
		}

		argV := reflect.New(field.Type)
		inter := argV.Interface()

		name := getNameInterface(inter)
		service := getPreparingServiceByName(c, name)

		st.Field(i).Set(reflect.ValueOf(service.instance))
	}
}

func injectByFunc(c *Container, s *service) {
	if !s.config.enableFunctionInject {
		return
	}

	var callArgs []reflect.Value

	serviceV := reflect.ValueOf(s.instance)

	funcV := serviceV.MethodByName(s.config.functionInjectName)
	funcT := funcV.Type()
	for i := 0; i < funcT.NumIn(); i++ {
		argT := funcT.In(i)
		argV := reflect.New(argT)
		inter := argV.Interface()
		name := getNameInterface(inter)

		service := getPreparingServiceByName(c, name)

		callArgs = append(callArgs, reflect.ValueOf(service.instance))
	}
	funcV.Call(callArgs)
}

func getName[T any]() string {
	inst := new(T)

	return getNameInterface(inst)
}
func getNameInterface(i any) string {
	return reflect.TypeOf(i).String()
}
