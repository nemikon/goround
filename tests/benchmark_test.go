package tests

import (
	"testing"

	"github.com/nemikon/goround"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type BService1 struct{}

func (b *BService1) InjectDependencies() {}

type BService2 struct {
	S1 *BService1
}

func (b *BService2) InjectDependencies(s1 *BService1) {
	b.S1 = s1
}

type BService3 struct {
	S2 *BService2
	S1 *BService1
}

func (b *BService3) InjectDependencies(s2 *BService2, s1 *BService1) {
	b.S1 = s1
	b.S2 = s2
}

type BService4 struct {
	S3 *BService3
	S2 *BService2
	S1 *BService1
}

func (b *BService4) InjectDependencies(s3 *BService3, s2 *BService2, s1 *BService1) {
	b.S1 = s1
	b.S2 = s2
	b.S3 = s3
}

type BService5 struct {
	S4 *BService4
	S3 *BService3
	S2 *BService2
	S1 *BService1
}

func (b *BService5) InjectDependencies(s4 *BService4, s3 *BService3, s2 *BService2, s1 *BService1) {
	b.S1 = s1
	b.S2 = s2
	b.S3 = s3
	b.S4 = s4
}

type BService6 struct {
	S5 *BService5
	S4 *BService4
	S3 *BService3
	S2 *BService2
	S1 *BService1
}

func (b *BService6) InjectDependencies(s5 *BService5, s4 *BService4, s3 *BService3, s2 *BService2, s1 *BService1) {
	b.S1 = s1
	b.S2 = s2
	b.S3 = s3
	b.S4 = s4
	b.S5 = s5
}

type BService7 struct {
	S6 *BService6
	S5 *BService5
	S4 *BService4
	S3 *BService3
	S2 *BService2
	S1 *BService1
}

func (b *BService7) InjectDependencies(s6 *BService6, s5 *BService5, s4 *BService4, s3 *BService3, s2 *BService2, s1 *BService1) {
	b.S1 = s1
	b.S2 = s2
	b.S3 = s3
	b.S4 = s4
	b.S5 = s5
	b.S6 = s6
}

type BService8 struct {
	S7 *BService7
	S6 *BService6
	S5 *BService5
	S4 *BService4
	S3 *BService3
	S2 *BService2
	S1 *BService1
}

func (b *BService8) InjectDependencies(s7 *BService7, s6 *BService6, s5 *BService5, s4 *BService4, s3 *BService3, s2 *BService2, s1 *BService1) {
	b.S1 = s1
	b.S2 = s2
	b.S3 = s3
	b.S4 = s4
	b.S5 = s5
	b.S6 = s6
	b.S7 = s7
}

type BService9 struct {
	S8 *BService8
	S7 *BService7
	S6 *BService6
	S5 *BService5
	S4 *BService4
	S3 *BService3
	S2 *BService2
	S1 *BService1
}

func (b *BService9) InjectDependencies(s8 *BService8, s7 *BService7, s6 *BService6, s5 *BService5, s4 *BService4, s3 *BService3, s2 *BService2, s1 *BService1) {
	b.S1 = s1
	b.S2 = s2
	b.S3 = s3
	b.S4 = s4
	b.S5 = s5
	b.S6 = s6
	b.S7 = s7
	b.S8 = s8
}

type BService10 struct {
	S9 *BService9
	S8 *BService8
	S7 *BService7
	S6 *BService6
	S5 *BService5
	S4 *BService4
	S3 *BService3
	S2 *BService2
	S1 *BService1
}

func (b *BService10) InjectDependencies(s9 *BService9, s8 *BService8, s7 *BService7, s6 *BService6, s5 *BService5, s4 *BService4, s3 *BService3, s2 *BService2, s1 *BService1) {
	b.S1 = s1
	b.S2 = s2
	b.S3 = s3
	b.S4 = s4
	b.S5 = s5
	b.S6 = s6
	b.S7 = s7
	b.S8 = s8
	b.S9 = s9
}

func TestBenchmarkFunctionInject(t *testing.T) {
	c := goround.NewContainer(goround.WithFunctionInject(true))
	registerBenchmarkServices(c)
	verifyBenchmarkServices(t, c)
}
func TestBenchmarkStructInject(t *testing.T) {
	c := goround.NewContainer(goround.WithStructInject(true), goround.WithStructDefaultInject(true))
	registerBenchmarkServices(c)
	verifyBenchmarkServices(t, c)
}
func BenchmarkFunctionInject(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		c := goround.NewContainer(goround.WithFunctionInject(true))
		registerBenchmarkServices(c)
		b.StartTimer()

		goround.Get[*BService10](c)
	}
}

func BenchmarkStructInject(b *testing.B) {
	for b.Loop() {
		b.StopTimer()
		c := goround.NewContainer(goround.WithStructInject(true), goround.WithStructDefaultInject(true))
		registerBenchmarkServices(c)
		b.StartTimer()

		goround.Get[*BService10](c)
	}
}

func verifyBenchmarkServices(t *testing.T, c *goround.Container) {
	s1 := goround.Get[*BService1](c)
	s2 := goround.Get[*BService2](c)
	s3 := goround.Get[*BService3](c)
	s4 := goround.Get[*BService4](c)
	s5 := goround.Get[*BService5](c)
	s6 := goround.Get[*BService6](c)
	s7 := goround.Get[*BService7](c)
	s8 := goround.Get[*BService8](c)
	s9 := goround.Get[*BService9](c)
	s10 := goround.Get[*BService10](c)

	require.NotNil(t, s1)
	require.NotNil(t, s2)
	require.NotNil(t, s3)
	require.NotNil(t, s4)
	require.NotNil(t, s5)
	require.NotNil(t, s6)
	require.NotNil(t, s7)
	require.NotNil(t, s8)
	require.NotNil(t, s9)
	require.NotNil(t, s10)

	assert.Same(t, s10.S1, s1)
	assert.Same(t, s10.S2, s2)
	assert.Same(t, s10.S3, s3)
	assert.Same(t, s10.S4, s4)
	assert.Same(t, s10.S5, s5)
	assert.Same(t, s10.S6, s6)
	assert.Same(t, s10.S7, s7)
	assert.Same(t, s10.S8, s8)
	assert.Same(t, s10.S9, s9)

	assert.Same(t, s9.S1, s1)
	assert.Same(t, s9.S2, s2)
	assert.Same(t, s9.S3, s3)
	assert.Same(t, s9.S4, s4)
	assert.Same(t, s9.S5, s5)
	assert.Same(t, s9.S6, s6)
	assert.Same(t, s9.S7, s7)
	assert.Same(t, s9.S8, s8)

	assert.Same(t, s8.S1, s1)
	assert.Same(t, s8.S2, s2)
	assert.Same(t, s8.S3, s3)
	assert.Same(t, s8.S4, s4)
	assert.Same(t, s8.S5, s5)
	assert.Same(t, s8.S6, s6)
	assert.Same(t, s8.S7, s7)

	assert.Same(t, s7.S1, s1)
	assert.Same(t, s7.S2, s2)
	assert.Same(t, s7.S3, s3)
	assert.Same(t, s7.S4, s4)
	assert.Same(t, s7.S5, s5)
	assert.Same(t, s7.S6, s6)

	assert.Same(t, s6.S1, s1)
	assert.Same(t, s6.S2, s2)
	assert.Same(t, s6.S3, s3)
	assert.Same(t, s6.S4, s4)
	assert.Same(t, s6.S5, s5)

	assert.Same(t, s5.S1, s1)
	assert.Same(t, s5.S2, s2)
	assert.Same(t, s5.S3, s3)
	assert.Same(t, s5.S4, s4)

	assert.Same(t, s4.S1, s1)
	assert.Same(t, s4.S2, s2)
	assert.Same(t, s4.S3, s3)

	assert.Same(t, s3.S1, s1)
	assert.Same(t, s3.S2, s2)

	assert.Same(t, s2.S1, s1)
}

func registerBenchmarkServices(c *goround.Container) {
	goround.Register(c, &BService1{})
	goround.Register(c, &BService2{})
	goround.Register(c, &BService3{})
	goround.Register(c, &BService4{})
	goround.Register(c, &BService5{})
	goround.Register(c, &BService6{})
	goround.Register(c, &BService7{})
	goround.Register(c, &BService8{})
	goround.Register(c, &BService9{})
	goround.Register(c, &BService10{})
}
