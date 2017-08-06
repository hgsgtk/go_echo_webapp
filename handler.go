package main

import(
  "net/http"
  "github.com/labstack/echo"
)

func setRoute(e *echo.Echo){
  // ルーティングハンドラ
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
