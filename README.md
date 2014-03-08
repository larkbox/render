# render [![wercker status](https://app.wercker.com/status/fcf6b26a1b41f53540200b1949b48dec "wercker status")](https://app.wercker.com/project/bykey/fcf6b26a1b41f53540200b1949b48dec)
Martini middleware/handler for easily rendering serialized JSON and HTML template responses.

[API Reference](http://godoc.org/github.com/martini-contrib/render)
## Hacked
对martini-contrib的render进行了改造，支持模板继承的方式。同时删除了yield和extensions。具体请看介绍：
## Usage
render uses Go's [html/template](http://golang.org/pkg/html/template/) package to render html templates.

###直接渲染hello.tmpl
~~~ go
// main.go
package main

import (
  "github.com/codegangsta/martini"
  "github.com/martini-contrib/render"
)

func main() {
  m := martini.Classic()
  // render html templates from templates directory
  m.Use(render.Renderer())

  m.Get("/", func(r render.Render) {
    r.HTML(200, "hello.tmpl", "jeremy")
  })

  m.Run()
}

~~~

~~~ html
<!-- templates/hello.tmpl -->
<h2>Hello {{.}}!</h2>
~~~

###layout.tmpl作为hello.tmpl的基础模板
~~~ go
// main.go
package main

import (
  "github.com/codegangsta/martini"
  "github.com/martini-contrib/render"
)

func main() {
  m := martini.Classic()
  // render html templates from templates directory
  m.Use(Renderer(Options{
	 Layout:    "layout.tmpl",
  }))
  

  m.Get("/", func(r render.Render) {
    r.HTML(200, "hello.tmpl", "jeremy")
  })

  m.Run()
}

~~~

~~~ html
<!-- templates/hello.tmpl -->
{{define "header"}}
	header
{{end}}

{{define "content"}}
	<h2>Hello {{.}}!</h2>
{{end}}

{{define "footer"}}
	footer
{{end}}
~~~
~~~ html
<!-- templates/layout.tmpl -->
{{template "header" .}}
{{template "content" .}}
{{template "footer" .}}
~~~


### Options
`render.Renderer` comes with a variety of configuration options:

~~~ go
// ...
m.Use(render.Renderer(render.Options{
  Directory: "templates", // Specify what path to load the templates from.
  Layout: "layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
  Funcs: []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
  Delims: render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
  Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
  IndentJSON: true, // Output human readable JSON
}))
// ...
~~~

### Cached template
生产环境下，将解析过的template缓存起来
~~~go
	if martini.Env == martini.Prod {
			r.tmpls[key] = t
		}
~~~

### Character Encodings
The `render.Renderer` middleware will automatically set the proper Content-Type header based on which function you call. See below for an example of what the default settings would output (note that UTF-8 is the default):
~~~ go
// main.go
package main

import (
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/render"
)

func main() {
  m := martini.Classic()
  m.Use(render.Renderer())

  // This will set the Content-Type header to "text/html; charset=UTF-8"
  m.Get("/", func(r render.Render) {
    r.HTML(200, "hello.tmpl", "world")
  })

  // This will set the Content-Type header to "application/json; charset=UTF-8"
  m.Get("/api", func(r render.Render) {
    r.JSON(200, map[string]interface{}{"hello": "world"})
  })

  m.Run()
}

~~~

In order to change the charset, you can set the `Charset` within the `render.Options` to your encoding value:
~~~ go
// main.go
package main

import (
  "github.com/codegangsta/martini"
  "github.com/codegangsta/martini-contrib/render"
)

func main() {
  m := martini.Classic()
  m.Use(render.Renderer(render.Options{
    Charset: "ISO-8859-1",
  }))

  // This is set the Content-Type to "text/html; charset=ISO-8859-1"
  m.Get("/", func(r render.Render) {
    r.HTML(200, "hello.tmpl", "world")
  })

  // This is set the Content-Type to "application/json; charset=ISO-8859-1"
  m.Get("/api", func(r render.Render) {
    r.JSON(200, map[string]interface{}{"hello": "world"})
  })

  m.Run()
}

~~~

## Authors
* [Jeremy Saenz](http://github.com/codegangsta)
* [Cory Jacobsen](http://github.com/cojac)
* [larkyang](http://github.com/rnoldo)
