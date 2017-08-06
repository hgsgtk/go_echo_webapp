package main

import(
  "html/template"
  "io"

  "github.com/labstack/echo"
)

// HTMLテンプレートを利用すためのRenderer Interface
type Template struct{
}

// HTMLテンプレートのデータを埋め込んだ結果をWriterに書き込む
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
  if t, ok := templates[name]; ok{
    return t.ExecuteTemplate(w, "layout.html", data)
  }
  c.Echo().Logger.Debugf("Template[%s] Not Found.", name)
  return templates["error"].ExecuteTemplate(w, "layout.html", "Internal Server Error")
  }

// 各HTMLテンプレートに共通レイアウトを適用した結果を保存する
func loadTemplates(){
  var baseTemplate = "templates/layout.html"
  templates = make(map[string]*template.Template)
  templates["index"]        = template.Must(
    template.ParseFiles(baseTemplate, "templates/index.html"))
  templates["error"]        = template.Must(
    template.ParseFiles(baseTemplate, "templates/error.html"))
  templates["user"]         = template.Must(
    template.ParseFiles(baseTemplate, "templates/user.html"))
  templates["login"]        = template.Must(
    template.ParseFiles(baseTemplate, "templates/login.html"))
  templates["admin"]        = template.Must(
    template.ParseFiles(baseTemplate, "templates/admin.html"))
  templates["admin_users"]  = template.Must(
    template.ParseFiles(baseTemplate, "templates/admin_users.html"))
}
