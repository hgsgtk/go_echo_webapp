package session

import(
  "errors"
  "time"

  "github.com/labstack/echo"
)

// セッションを一意に識別するID
type ID string

// セッションデータと整合性トークンを保持する構造体
type Store struct{
  Data              map[string]string
  ConsistencyToken  string
}

// Sessionの操作・管理するManager
type Manager struct{
  stopCh    chan  struct{}
  commandCh chan  command
  stopGCCh  chan  struct{}
}

// Managerの開始
func(m *Manager)Start(echo *echo.Echo){
  e = echo
  go m.mainLoop()
  time.Sleep(100 * time.Millisecond)
  go m.gcLoop()
}

// Managerの停止
func(m *Manager)Stop(){
  m.stopGCCh <- struct{}{}
  time.Sleep(100 * time.Millisecond)
  m.stopCh <- struct{}{}
}

// Sessionの作成
func(m *Manager)Create()(ID, error){
  respCh := make(chan response, 1)
  defer close(respCh)
  cmd := command{commandCreate, nil, respCh}
  m.commandCh <- cmd
  resp := <-respCh
  var res ID
  if resp.err != nil{
    e.Logger.Debugf("Session Create Error. [%s]", resp.err)
    return res, resp.err
  }
  if res, ok := resp.result[0].(ID); ok{
    return res, nil
  }
  e.Logger.Debugf("Session Create Error. [%s]", ErrorOther)
  return res, ErrorOther
}

// Datastoreの読み出し
func(m *Manager)LoadStore(sessionID ID)(Store, error){
  respCh := make(chan response, 1)
  defer close(respCh)
  req := []interface{}{sessionID}
  cmd := command{commandLoadStore, req, respCh}
  m.commandCh <- cmd
  resp := <-respCh
  var res Store
  if resp.err != nil{
    e.Logger.Debugf("Session[%s] Load store Error. [%s]", sessionID, resp.err)
    return res, resp.err
  }
  if res, ok := resp.result[0].(Store); ok{
    return res, nil
  }
  e.Logger.Debugf("Session[%s] Load store Error. [%s]", sessionID, ErrorOther)
  return res, ErrorOther
}

// Datastoreの保存
func(m *Manager)SaveStore(sessionID ID, sessionStore Store) error{
  respCh := make(chan response, 1)
  defer close(respCh)
  req := []interface{}{sessionID, sessionStore}
  cmd := command{commandSaveStore, req, respCh}
  m.commandCh <- cmd
  resp := <-respCh
  if resp.err != nil{
    e.Logger.Debugf("Session[%s] Save store Error. [%s]", sessionID, resp.err)
    return resp.err
  }
  return nil
}

// Sessionの削除
func(m *Manager)Delete(sessionID ID) error{
  respCh := make(chan response, 1)
  defer close(respCh)
  req := []interface{}{sessionID}
  cmd := command{commandDelete, req, respCh}
  m.commandCh <- cmd
  resp := <-respCh
  if resp.err != nil{
    e.Logger.Debugf("Session[%s] Delete Error. [%s]", sessionID, resp.err)
    return resp.err
  }
  return nil
}

// 期限切れSessionの削除
func(m *Manager)DeleteExpired() error{
  respCh := make(chan response, 1)
  defer close(respCh)
  cmd := command{commandDelete, nil, respCh}
  m.commandCh <- cmd
  resp := <-respCh
  if resp.err != nil{
    e.Logger.Debugf("Session DeleteExpired Error. [%s]", resp.err)
  }
  return nil
}

// Managerが返す各エラーのinstance
var(
  ErrorBadParameter     = errors.New("Bad Parameter")
  ErrorNotFound         = errors.New("Not Found")
  ErrorInvalidToken     = errors.New("Invalid Token")
  ErrorInvalidCommand   = errors.New("Invalid Command")
  ErrorNotImplemented   = errors.New("Not Implemented")
  ErrorOther            = errors.New("Other")
)
