package main

import(
  "net/http"
  "./model"
  "github.com/labstack/echo"
)

func setRoute(e *echo.Echo){
  e.GET("/", getIndex)
  e.GET("/login", getLogin)
  e.POST("/login", postLogin)
  e.POST("/logout", postLogout)
  e.GET("/users/:user_id", handleUser)
  e.POST("/users/:user_id", handleUser)

  // 管理者のみ参照のページ
  admin := e.Group("/admin", MiddlewareAuthAdmin)
  admin.GET("", handleAdminTop)
  admin.POST("", handleAdminTop)
  admin.GET("/users", getAdminUser)

}

// GET:/
func getIndex(c echo.Context) error{
  return c.Render(http.StatusOK, "index", "Hello, World!")

}

// GET:/users/:user_id
// POST:/users/:user_id
func handleUser(c echo.Context)error{
  userID := c.Param("user_id")
  err := CheckUserID(c, userID)
  if err != nil{
    c.Echo().Logger.Debugf("User Page[%s] Role Error. [%s]", userID, err)
    msg := "ログインしていません。"
    return c.Render(http.StatusOK, "error", msg)
  }
  users, err := userDA.FindByUserID(userID, model.FindFirst)
  if err != nil{
    return c.Render(http.StatusOK, "error", err)
  }
  user := users[0]
  return c.Render(http.StatusOK, "user", user)
}

// GET:/admin
// POST:/admin
func handleAdminTop(c echo.Context) error {
  return c.Render(http.StatusOK, "admin", nil)
}

// GET:/admin/users
func getAdminUser(c echo.Context)error{
  users, err := userDA.FindAll()
  if err != nil{
    return c.Render(http.StatusOK, "errors", err)
  }
  return c.Render(http.StatusOK, "admin_users", users)
}

// GET:/Login
func getLogin(c echo.Context) error {
  return c.Render(http.StatusOK, "login", nil)
}

// POST:/Login
func postLogin(c echo.Context) error{
  userID := c.FormValue("userid")
  password := c.FormValue("password")
  err := UserLogin(c, userID, password)
  if err != nil{
    c.Echo().Logger.Debugf("User[%s] Login Error. [%s]", userID, err)
    msg := "ユーザーIDまたはパスワードが誤っています。"
    data := map[string]string{"user_id": userID, "password": "", "msg": msg}
    return c.Render(http.StatusOK, "login", data)
  }
  // ログインしたユーザーが管理者かチェック
  isAdmin, err := CheckRoleByUserID(userID, model.RoleAdmin)
  if err != nil{
    c.Echo().Logger.Debugf("Admin Role Check Error. [%s]", userID, err)
    isAdmin = false
  }
  if isAdmin{
    c.Echo().Logger.Debugf("User is Admin. [%s]", userID)
    return c.Redirect(http.StatusTemporaryRedirect, "/admin")
  }
  return c.Redirect(http.StatusTemporaryRedirect, "/users"+userID)
}

// POST:/logout
func postLogout(c echo.Context) error {
  err := UserLogout(c)
  if err != nil{
    c.Echo().Logger.Debugf("User Logout Error. [%s]", err)
    return c.Render(http.StatusOK, "login", nil)
  }
  msg := "ログアウトしました。"
  data := map[string]string{"user_id": "", "password": "", "msg": msg}
  return c.Render(http.StatusOK, "login", data)
}
