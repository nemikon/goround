package tests

import (
	"testing"
	"time"

	"github.com/nemikon/goround"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type IService1 interface {
	IAmService1()
}

type IService2 interface {
	IAmService2()
}

type FService1 struct {
	s2 IService2
}

func (s *FService1) IAmService1() {}
func (s *FService1) InjectDependencies(s2 IService2) {
	s.s2 = s2
}

type FService2 struct {
	s1 IService1
}

func (s *FService2) IAmService2() {}
func (s *FService2) InjectDependencies(s1 IService1) {
	s.s1 = s1
}

type SService1 struct {
	S2 IService2 `goround:"inject"`
}

func (s *SService1) IAmService1() {}

type SService2 struct {
	S1 IService1 `goround:"inject"`
}

func (s *SService2) IAmService2() {}

func TestProvideBasic(t *testing.T) {
	t.Parallel()

	c := goround.NewContainer(goround.WithFunctionInject(true), goround.WithStructInject(true))

	type DB struct{ Con string }
	db := DB{"connectionString"}
	dbPtr := &DB{"ptrConnectionString"}

	goround.RegisterBasicValue(c, db)
	goround.RegisterBasicValue(c, dbPtr)

	outDb := goround.Get[DB](c)
	outDbPtr := goround.Get[*DB](c)

	assert.Equal(t, db, outDb)
	assert.Equal(t, dbPtr, outDbPtr)
	assert.Same(t, dbPtr, outDbPtr)
}

func TestStructInject(t *testing.T) {
	t.Parallel()

	c := goround.NewContainer(goround.WithStructInject(true))

	emptyS1 := &SService1{}
	emptyS2 := &SService2{}

	goround.Register[IService1](c, emptyS1)
	goround.Register[IService2](c, emptyS2)

	s1 := goround.Get[IService1](c)
	s2 := goround.Get[IService2](c)

	require.Same(t, emptyS1, s1)
	require.Same(t, emptyS2, s2)

	require.Same(t, emptyS2, emptyS1.S2)
	require.Same(t, emptyS1, emptyS2.S1)

}

func TestFunctionInject(t *testing.T) {
	t.Parallel()

	c := goround.NewContainer(goround.WithFunctionInject(true))
	emptyS1 := &FService1{}
	emptyS2 := &FService2{}

	goround.Register[IService1](c, emptyS1)
	goround.Register[IService2](c, emptyS2)

	s1 := goround.Get[IService1](c)
	s2 := goround.Get[IService2](c)

	require.Same(t, emptyS1, s1)
	require.Same(t, emptyS2, s2)

	require.Same(t, emptyS2, emptyS1.s2)
	require.Same(t, emptyS1, emptyS2.s1)
}

func TestDirectCircle(t *testing.T) {
	type Service struct {
		Dep *Service
	}

	c := goround.NewContainer(goround.WithStructInject(true), goround.WithStructDefaultInject(true))

	service := &Service{}
	goround.Register(c, service)

	outService := goround.Get[*Service](c)
	require.Same(t, service, outService)
	require.Same(t, service, outService.Dep)
}

type ConcurrentTestService struct {
	continueChan chan any
	dep          ConcurrentTestDependencie
}
type ConcurrentTestDependencie struct {
	content string
}

func (s *ConcurrentTestService) InjectDependencies(c ConcurrentTestDependencie) {
	<-s.continueChan
	s.dep = c
}

func TestConcurrentGet(t *testing.T) {
	t.Parallel()

	c := goround.NewContainer(goround.WithFunctionInject(true))

	s := &ConcurrentTestService{
		continueChan: make(chan any),
	}
	dep := ConcurrentTestDependencie{"hello world"}
	cont := func() {
		select {
		case s.continueChan <- 1:
		default:
		}
	}
	goround.Register(c, s)
	goround.RegisterBasicValue(c, dep)

	go func() {
		time.Sleep(200 * time.Millisecond)
		cont()
	}()

	go func() {
		outS := goround.Get[*ConcurrentTestService](c)
		assert.Equal(t, dep, outS.dep)
		cont()
	}()

	outS := goround.Get[*ConcurrentTestService](c)
	assert.Equal(t, dep, outS.dep)
	cont()
}
