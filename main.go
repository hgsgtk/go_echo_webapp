package main

import(
  "context"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"
  "io"
  "strconv"
  "html/template"
  "log"
  "fmt"

  "./session"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct{
  Name string `json:"name" xml:"name" form:"name" query:"name"`
  Email string `json:"email" xml:"email" form:"email" query:"email"`
}

type Customer struct{
  Id  int `json:id`
  Name  string  `json:name`
  Sex int `json:sex`
  Tel string  `json:tel`
}

var templates map[string]*template.Template

var sessionManager *session.Manager

type Template struct{
}

// セッションのCookieに関する設定
const(
  sessionCookieName string = "apisrv_session_id"
  sessionCookieExpire time.Duration = (1 * time.Hour)
)
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
  return templates[name].ExecuteTemplate(w, "layout.html", data)
}

func main(){
  // Create Echo Instance
  e := echo.New()

  // Renderer Setting to use templates
  t := &Template{}
  e.Renderer = t

  // Middleware Setting
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // File Path Setting
  e.Static("/public/css/", "./public/css/")
  e.Static("/public/js/", "./public/js/")
  e.Static("/public/img/", "./public/img/")


  // Routing Handleer
  e.GET("/", home)
  e.GET("/users", users)
  e.GET("/users/:id", showUser)
  e.GET("users/search", searchUser)
  e.POST("/save", save)
  e.POST("/users/save", saveUser)
  e.POST("/users", saveUsers)
  e.GET("/session_form", showSessionForm)
  e.GET("/session", showSession)
  e.POST("/session", postSession)
  e.POST("/session_delete", deleteSession)

  // セッション管理を開始
  sessionManager = &session.Manager{}
  sessionManager.Start(e)

  // サーバを開始
  e.Logger.Fatal(e.Start(":1323"))

  // 中断を検知したらリクエストの完了を10秒まで待ってサーバを終了する
  quit := make(chan os.Signal)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  <-quit
  ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
  defer cancel()
  if err := e.Shutdown(ctx); err != nil{
    e.Logger.Info(err)
    e.Close()
  }

  // セッション管理を停止
  sessionManager.Stop()

  // 終了ログが出るまで少し待つ
  time.Sleep(1 * time.Second)
}

func init(){
  loadTemplates()
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


func home(c echo.Context) error{
  return c.Render(http.StatusOK, "index", "Hello, World!")

}

func users(c echo.Context) error{
  db := connectDB()

  customerEx := Customer{}
  db.Find(&customerEx)
  log.Print(customerEx)

  data := struct{
    id int
    name string
    sex int
  }{
    customerEx.Id,
    customerEx.Name,
    customerEx.Sex,
  }
  return c.Render(http.StatusOK, "customers.html", data)
}


func showUser(c echo.Context) error{
  db := connectDB()

  customerEx := Customer{}
  paramId, _ := strconv.Atoi(c.Param("id"))
  customerEx.Id = paramId
  db.First(&customerEx)
  return c.String(http.StatusOK, customerEx.Name)
}

func searchUser(c echo.Context) error{
  var keyword string
  keyword = c.Param("keyword")
  var Info string
  Info = searchDB(keyword)

  return c.Render(http.StatusOK, "search.html", Info)
}

func save(c echo.Context) error{
  name := c.FormValue("name")
  email := c.FormValue("email")
  return c.String(http.StatusOK, "name: " + name +  ", email: " + email)
}

func saveUser(c echo.Context) error{
  name := c.FormValue("name")
  avatar, err := c.FormFile("avatar")
  if err != nil{
    return err
  }

  src, err := avatar.Open()
  if err != nil{
    return err
  }
  defer src.Close()

  dst, err := os.Create(avatar.Filename)
  if err != nil{
    return err
  }
  defer dst.Close()

  if _, err = io.Copy(dst,src); err != nil{
    return err
  }

  return c.HTML(http.StatusOK, "<b>Thank you! " + name + "</b>")
}

func saveUsers(c echo.Context) error{
  u := new(User)
  if err := c.Bind(u); err != nil{
    return err
  }
  return c.JSON(http.StatusCreated, u)
}

func showSessionForm(c echo.Context) error{
  return c.Render(http.StatusOK, "session_form", nil)
}

func showSession(c echo.Context) error{
  sessionID, err := readSessionID(c)
  if err != nil{
    return c.Render(http.StatusOK, "error", "Session cookie Read Error: "+err.Error())
  }
  sessionStore, err := sessionManager.LoadStore(sessionID)
  if err != nil{
    return c.Render(http.StatusOK, "error", "Session store Load Error: "+err.Error())
  }
  return c.Render(http.StatusOK, "session",
    map[string]interface{}{"msg": "", "data": sessionStore.Data})
}

func postSession(c echo.Context) error{
  key1 := c.FormValue("key1")
  key2 := c.FormValue("key2")
  sessionID, err := sessionManager.Create()
  if err != nil{
    return c.Render(http.StatusOK, "error", "Session Create Error: "+err.Error())
  }
  writeSessionID(c, sessionID)
  sessionStore, err := sessionManager.LoadStore(sessionID)
  if err != nil{
    return c.Render(http.StatusOK, "error", "Session store Load Error: "+err.Error())
  }
  sessionData := map[string]string{
    "key1": key1,
    "key2": key2,
  }
  sessionStore.Data = sessionData
  err = sessionManager.SaveStore(sessionID, sessionStore)
  if err != nil{
    return c.Render(http.StatusOK, "error", "Session store Save Error: "+err.Error())
  }
  return c.Render(http.StatusOK, "session_form",
    map[string]interface{}{"msg": "セッションデータを保存しました", "data": sessionStore.Data})
}

func deleteSession(c echo.Context) error{
  sessionID, err := readSessionID(c)
  if err != nil{
    return c.Render(http.StatusOK, "error", "Session cookie Read Error: "+err.Error())
  }
  err = sessionManager.Delete(sessionID)
  if err != nil{
    return c.Render(http.StatusOK, "error", "Session delete Error: "+err.Error())
  }
  return c.Render(http.StatusOK, "session",
    map[string]interface{}{"msg": "セッション" + sessionID + "を削除しました。"})
}

// ブラウザのcookieにsession.IDを書き込む
func writeSessionID(c echo.Context, sessionID session.ID) error{
  cookie := new(http.Cookie)
  cookie.Name = sessionCookieName
  cookie.Value = string(sessionID)
  cookie.Expires = time.Now().Add(sessionCookieExpire)
  c.SetCookie(cookie)
  return nil
}

// ブラウザのcookieからsession.IDを読み込む
func readSessionID(c echo.Context)(session.ID, error){
  var sessionID session.ID
  cookie, err := c.Cookie(sessionCookieName)
  if err != nil{
    return sessionID, err
  }
  sessionID = session.ID(cookie.Value)
  return sessionID, nil
}

func connectDB() *gorm.DB{
  DBMS := "mysql"
  USER := "root"
  PASS := "verysecret"
  PROTOCOL := "tcp(127.0.0.1:13306)"
  DBNAME := "customer"

  CONNECT := USER+":"+PASS+"@"+PROTOCOL+"/"+DBNAME
  db, err := gorm.Open(DBMS, CONNECT)

  if err != nil{
    panic(err.Error())
  }
  return db
}

func searchDB(keyword string) string{
  db := connectDB()

  var customer Customer
  db.First(&customer, "Name = ?", keyword)

  var Info string
  Info += "ID: " + fmt.Sprint(customer.Id)
  Info += "Name: " + fmt.Sprint(customer.Name)

  return Info
}
