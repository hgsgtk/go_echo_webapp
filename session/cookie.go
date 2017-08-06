package session

import(
  "net/http"
  "time"

  "../setting"
  "github.com/labstack/echo"
)

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
