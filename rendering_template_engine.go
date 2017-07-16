package main

import(
  "html/template"
  "io"
  "log"
  "net/http"

  "github.com/labstack/echo"
)

type TemplateRenderer struct{
  templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
  if viewContext, isMap := data.(map[string]interface{}); isMap{
    viewContext["reverse"] = c.Echo().Reverse
  }

  return t.templates.ExecuteTemplate(w, name, data)
}

func main(){
  e := echo.New()
  renderer := &TemplateRenderer{
    templates: template.Must(template.ParseGlob("*.html")),
  }
  e.Renderer = renderer

  e.GET("/smt", func(c echo.Context) error{
    return c.Render(http.StatusOK, "template.html", map[string]interface{}{
      "name": "Kazuki",
    })
  }).Name = "foobar"

  log.Fatal(e.Start(":8000"))
}
