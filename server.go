package main

import(
  "net/http"
  "github.com/labstack/echo"
  "os"
  "io"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/mysql"
  "strconv"
  "html/template"
  "log"
  "fmt"
)

type TemplateRenderer struct{
  templates *template.Template
}

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


func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
  if viewContext, isMap := data.(map[string]interface{}); isMap{
    viewContext["reverse"] = c.Echo().Reverse
  }
  return t.templates.ExecuteTemplate(w, name, data)
}


func main(){
  e := echo.New()
  renderer := &TemplateRenderer{
    templates: template.Must(template.ParseGlob("html/*.html")),
  }
  e.Renderer = renderer

  // Routing
  e.GET("/", home)
  e.GET("/users", users)
  e.GET("/users/:id", showUser)
  e.GET("users/search", searchUser)
  e.POST("/save", save)
  e.POST("/users/save", saveUser)
  e.POST("/users", saveUsers)
  e.Logger.Fatal(e.Start(":1323"))
}

func home(c echo.Context) error{
  return c.String(http.StatusOK, "Hello, World!")

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
