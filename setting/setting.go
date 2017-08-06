package setting

import(
  "time"
)

// サーバの動作に関する設定
var Server = server{}

type server struct{
  Port  string
}

// セッションに関する設定
var Session = session{}

type session struct{
  CookieName    string
  CookieExpire  time.Duration
}

func Load(){
  Server.Port = ":3000"
  Session.CookieName = "session_id"
  Session.CookieExpire = (1 * time.Hour)
}
