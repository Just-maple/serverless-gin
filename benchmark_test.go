package svrlessgin

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

type mockWriter struct {
	headers http.Header
}

func newMockWriter() *mockWriter {
	return &mockWriter{
		http.Header{},
	}
}

func (m *mockWriter) Header() (h http.Header) {
	return m.headers
}

func (m *mockWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockWriter) WriteHeader(int) {}

func runTestRequest(t *testing.T, r *gin.Engine, method, path string) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	w := newMockWriter()
	r.ServeHTTP(w, req)
}

func runRequest(B *testing.B, r *gin.Engine, method, path string) {
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		panic(err)
	}
	B.ReportAllocs()
	B.ResetTimer()
	//B.RunParallel(func(pb *testing.PB) {
	//	for pb.Next() {
	//		r.ServeHTTP(newMockWriter(), req)
	//	}
	//})
	w := newMockWriter()
	for i := 0; i < B.N; i++ {
		r.ServeHTTP(w, req)
	}
}

type ST struct {
	I int
	S string
}

type Ctl struct {
}

func (c2 Ctl) Response(c *gin.Context, ret interface{}, err error) {
	return
}

func (c2 Ctl) ParamHandler(c *gin.Context, params []interface{}) bool {
	return true
}

var _ GinIOController = &Ctl{}

func BenchmarkRun(b *testing.B) {
	e := gin.New()
	var f = NewWithController(&Ctl{})(func(ctx context.Context, param ST) (err error) {
		return nil
	})
	e.GET("/", f)
	runRequest(b, e, "GET", "/")
}

type service struct{}

func (s service) Func(ctx context.Context, param ST) (err error) {
	f, _ := GetServiceFunc(ctx)
	fmt.Printf("%v\n", f)
	return nil
}

func (s service) Func2(ctx context.Context, param ST) (err error) {
	f, _ := GetServiceFunc(ctx)
	fmt.Printf("%v\n", f)
	return nil
}

func TestServiceName(t *testing.T) {
	e := gin.New()
	f := NewWithController(&Ctl{})
	SetServiceFuncInject(true)
	e.GET("/", f(service{}.Func))
	e.GET("/2", f(service{}.Func2))
	runTestRequest(t, e, "GET", "/")
	runTestRequest(t, e, "GET", "/2")
}

func BenchmarkRunRaw(b *testing.B) {
	e := gin.New()
	var f = NewWithController(&Ctl{})(func() (err error) {
		return nil
	})
	e.GET("", f)
	runRequest(b, e, "GET", "")
}

func BenchmarkRunDef(b *testing.B) {
	e := gin.New()
	e.GET("", func(c *gin.Context) {
		var param ST
		_ = c.Bind(&param)
	})
	runRequest(b, e, "GET", "")
}

func BenchmarkRunDefRaw(b *testing.B) {
	e := gin.New()
	e.GET("", func(c *gin.Context) {
	})
	runRequest(b, e, "GET", "")
}
