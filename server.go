package main

import(
  "context"
  "html/template"
  "os"
  "os/signal"
  "syscall"
  "time"
  "log"

  "./session"
  "./setting"
  "github.com/labstack/echo"
  "github.com/labstack/echo/middleware"
)

var templates map[string]*template.Template

var sessionManager *session.Manager

func main(){
  // Create Echo Instance
  e := echo.New()

  // Renderer Setting to use templates
  t := &Template{}
  e.Renderer = t

  // Middleware Setting
  e.Use(middleware.Logger())
  e.Use(middleware.Recover())

  // 静的ファイルを配置するルーティングを設定
  setStaticRoute(e)

  // 各ルーティングに対するハンドルを設定
  setRoute(e)

  // セッション管理を開始
  sessionManager = &session.Manager{}
  sessionManager.Start(e)

  // サーバ開始
  go func(){
    if err := e.Start(setting.Server.Port); err != nil{
      e.Logger.Info("shutting down the server")
    }
  }()

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

// 初期化を行う
func init(){
  // 設定の読み込み
  setting.Load()
  // HTMLテンプレートの読み込み
  loadTemplates()
}