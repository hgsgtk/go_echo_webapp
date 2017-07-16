package main

import(
  "net/http"
  "github.com/labstack/echo"
  "os"
  "io"
)

type User struct{
  Name string `json:"name" xml:"name" form:"name" query:"name"`
  Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func main(){
  e := echo.New()

  // Routing
  e.GET("/", home)
  e.GET("/users/:id", getUser)
  e.GET("/show", show)
  e.POST("/save", save)
  e.POST("/users/save", saveUser)
  e.POST("/users", saveUsers)
  e.Logger.Fatal(e.Start(":1323"))
}

func home(c echo.Context) error{
  return c.String(http.StatusOK, "Hello, World!")

}
func getUser(c echo.Context) error{
  id := c.Param("id")
  return c.String(http.StatusOK, id)
}

func show(c echo.Context) error{
  team := c.QueryParam("team")
  member := c.QueryParam("member")
  return c.String(http.StatusOK, "team: " + team + ", member: " + member)
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
