package main

import(
  "html/template"
  "io"

  "github.com/labstack/echo"
)

// HTMLテンプレートを利用するためのレンダリングインターフェース
type Template struct{
}

// HTMLテンプレートのデータを埋め込んだ結果をWriterに書き込む
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
  return templates[name].ExecuteTemplate(w, "layout.html", data)
}

// 各HTMLテンプレートに共通レイアウトを適用した結果を保存する
func loadTemplates(){
  var baseTemplate = "templates/layout.html"
  templates = make(map[string]*template.Template)
  templates["index"]        = template.Must(
    template.ParseFiles(baseTemplate, "templates/hello.html"))
  templates["session_form"] = template.Must(
    template.ParseFiles(baseTemplate, "templates/session_form.html"))
  templates["session"]      = template.Must(
    template.ParseFiles(baseTemplate, "templates/session.html"))
  templates["error"]        = template.Must(
    template.ParseFiles(baseTemplate, "templates/error.html"))
}
