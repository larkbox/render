package render

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/codegangsta/martini"
)

type Greeting struct {
	One string `json:"one"`
	Two string `json:"two"`
}

func Test_Render_JSON(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
	// nothing here to configure
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.JSON(300, Greeting{"hello", "world"})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 300)
	expect(t, res.Header().Get(ContentType), ContentJSON+"; charset=UTF-8")
	expect(t, res.Body.String(), `{"one":"hello","two":"world"}`)
}

func Test_Render_Indented_JSON(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		IndentJSON: true,
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.JSON(300, Greeting{"hello", "world"})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 300)
	expect(t, res.Header().Get(ContentType), ContentJSON+"; charset=UTF-8")
	expect(t, res.Body.String(), `{
  "one": "hello",
  "two": "world"
}`)
}

// func Test_Render_Bad_HTML(t *testing.T) {
// 	m := martini.Classic()
// 	m.Use(Renderer(Options{
// 		Directory: "fixtures/basic",
// 	}))

// 	// routing
// 	m.Get("/foobar", func(r Render) {
// 		r.HTML(200, "nope.tmpl", nil)
// 	})

// 	res := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/foobar", nil)

// 	m.ServeHTTP(res, req)

// 	expect(t, res.Code, 500)
// 	// expect(t, res.Body.String(), "html/template: \"nope\" is undefined\n")
// }

func Test_Render_NO_LAYOUT_HTML(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "hello.tmpl", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "<h1>Hello jeremy</h1>\n")
}

func Test_Render_Funcs(t *testing.T) {

	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/custom_funcs",
		Funcs: []template.FuncMap{
			{
				"myCustomFunc": func() string {
					return "My custom function"
				},
			},
		},
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "index.tmpl", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Body.String(), "My custom function\n")
}

func Test_Render_Layout(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
		Layout:    "layout.tmpl",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "content.tmpl", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Body.String(), "\n\theader\n\n\tcontent\n\n\tfooter\n")
}

func Test_Render_Nested_HTML(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
		Layout:    "layout.tmpl",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "admin/index.tmpl", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "\n\theader-admin\n\n\tcontent-admin\n\n\tfooter-admin\n")
}

func Test_Render_Delimiters(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Delims:    Delims{"{[{", "}]}"},
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "delims.tmpl", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "<h1>Hello jeremy</h1>")
}

func Test_Render_Error404(t *testing.T) {
	res := httptest.NewRecorder()
	r := renderer{res, nil, nil, Options{}, ""}
	r.Error(404)
	expect(t, res.Code, 404)
}

func Test_Render_Error500(t *testing.T) {
	res := httptest.NewRecorder()
	r := renderer{res, nil, nil, Options{}, ""}
	r.Error(500)
	expect(t, res.Code, 500)
}

func Test_Render_Redirect_Default(t *testing.T) {
	url, _ := url.Parse("http://localhost/path/one")
	req := http.Request{
		Method: "GET",
		URL:    url,
	}
	res := httptest.NewRecorder()

	r := renderer{res, &req, nil, Options{}, ""}
	r.Redirect("two")

	expect(t, res.Code, 302)
	expect(t, res.HeaderMap["Location"][0], "/path/two")
}

func Test_Render_Redirect_Code(t *testing.T) {
	url, _ := url.Parse("http://localhost/path/one")
	req := http.Request{
		Method: "GET",
		URL:    url,
	}
	res := httptest.NewRecorder()

	r := renderer{res, &req, nil, Options{}, ""}
	r.Redirect("two", 307)

	expect(t, res.Code, 307)
	expect(t, res.HeaderMap["Location"][0], "/path/two")
}

func Test_Render_Charset_JSON(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Charset: "foobar",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.JSON(300, Greeting{"hello", "world"})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 300)
	expect(t, res.Header().Get(ContentType), ContentJSON+"; charset=foobar")
	expect(t, res.Body.String(), `{"one":"hello","two":"world"}`)
}

func Test_Render_Default_Charset_HTML(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "hello.tmpl", "jeremy")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	// ContentLength should be deferred to the ResponseWriter and not Render
	expect(t, res.Header().Get(ContentLength), "")
	expect(t, res.Body.String(), "<h1>Hello jeremy</h1>\n")
}

func Test_Render_Override_Layout(t *testing.T) {
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
		Layout:    "layout.tmpl",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "content.tmpl", "jeremy", HTMLOptions{
			Layout: "another_layout.tmpl",
		})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "another\n\theader\n\nanother\n\tcontent\n\nanother\n\tfooter\n\n")
}

func Test_Render_NoRace(t *testing.T) {
	// This test used to fail if run with -race
	m := martini.Classic()
	m.Use(Renderer(Options{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r Render) {
		r.HTML(200, "hello.tmpl", "world")
	})

	done := make(chan bool)
	doreq := func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/foobar", nil)

		m.ServeHTTP(res, req)

		expect(t, res.Code, 200)
		expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
		// ContentLength should be deferred to the ResponseWriter and not Render
		expect(t, res.Header().Get(ContentLength), "")
		expect(t, res.Body.String(), "<h1>Hello world</h1>\n")
		done <- true
	}
	// Run two requests to check there is no race condition
	go doreq()
	go doreq()
	<-done
	<-done
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
