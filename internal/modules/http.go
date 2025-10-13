package modules

import (
	"bytes"
	"encoding/json"
  "fmt"
	"io"
	netHttp "net/http"
  "os"
  "strings"
  "sync"

	"github.com/dop251/goja"
  "github.com/gorilla/websocket"

	"github.com/douglasjordan2/dougless/internal/event"
)

type HTTP struct {
  vm        *goja.Runtime
  eventLoop *event.Loop
}

func NewHTTP(eventLoop *event.Loop) *HTTP {
  return &HTTP{
    eventLoop: eventLoop,
  }
}

func (http *HTTP) Export(vm *goja.Runtime) goja.Value {
  http.vm = vm
  obj := vm.NewObject()

  obj.Set("get", http.get)
  obj.Set("post", http.post)
  obj.Set("createServer", http.createServer)

  return obj
}

func (http *HTTP) get(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 2 {
    panic(http.vm.ToValue("GET requires a URL and a callback"))
  }

  url := call.Arguments[0].String()
  callback, ok := goja.AssertFunction(call.Arguments[1])
  if !ok {
    panic(http.vm.ToValue("second argument must be a callback function"))
  }

  http.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
      resp, err := netHttp.Get(url)

      var errArg, dataArg goja.Value
      if err != nil {
        errArg = http.vm.ToValue(err.Error())
        dataArg = goja.Undefined()
      } else {
        defer resp.Body.Close()
        body, readErr := io.ReadAll(resp.Body)

        if readErr != nil {
          errArg = http.vm.ToValue(readErr.Error())
          dataArg = goja.Undefined()
        } else {
          responseObj := http.vm.NewObject()
          responseObj.Set("status", resp.Status)
          responseObj.Set("statusCode", resp.StatusCode)
          responseObj.Set("body", string(body))

          headersObj := http.vm.NewObject()
          for key, values := range resp.Header {
            if len(values) == 1 {
              headersObj.Set(key, values[0])
            } else if len(values) > 1 {
              headersObj.Set(key, values)
            }
          }
          responseObj.Set("headers", headersObj)

          errArg = goja.Null()
          dataArg = responseObj
        }
      }

      callback(goja.Undefined(), errArg, dataArg)
    },
  })

  return goja.Undefined()
}

func (http *HTTP) post(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 3 {
    panic(http.vm.ToValue("POST requires a URL, payload, and a callback"))
  }

  url := call.Arguments[0].String()
  payload := call.Arguments[1].Export()
  callback, ok := goja.AssertFunction(call.Arguments[2])
  if !ok {
    panic(http.vm.ToValue("last argument must be a callback function"))
  }

  contentType := "application/json"
  dataMap, isMap := payload.(map[string]any)

  if isMap {
    if ct, exists := dataMap["contentType"]; exists {
      contentType = ct.(string)
      delete(dataMap, "contentType")
      payload = dataMap
    }
  }

  http.eventLoop.ScheduleTask(&event.Task{
    Callback: func() {
      jsonBytes, _ := json.Marshal(payload)
      body := bytes.NewBuffer(jsonBytes)

      resp, err := netHttp.Post(url, contentType, body)

      var errArg, dataArg goja.Value
      if err != nil {
        errArg = http.vm.ToValue(err.Error())
        dataArg = goja.Undefined()
      } else {
        defer resp.Body.Close()
        body, readErr := io.ReadAll(resp.Body)

        if readErr != nil {
          errArg = http.vm.ToValue(readErr.Error())
          dataArg = goja.Undefined()
        } else {
          responseObj := http.vm.NewObject()
          responseObj.Set("status", resp.Status)
          responseObj.Set("statusCode", resp.StatusCode)
          responseObj.Set("body", string(body))

          headersObj := http.vm.NewObject()
          for key, values := range resp.Header {
            if len(values) == 1 {
              headersObj.Set(key, values[0])
            } else if len(values) > 1 {
              headersObj.Set(key, values)
            }
          }
          responseObj.Set("headers", headersObj)

          errArg = goja.Null()
          dataArg = responseObj
        }
      }

      callback(goja.Undefined(), errArg, dataArg)
    },
  })

  return goja.Undefined()
}

func (http *HTTP) createRequestObject(r *netHttp.Request) goja.Value {
  reqObj := http.vm.NewObject()

  reqObj.Set("method", r.Method)
  reqObj.Set("url", r.URL.String())

  defer r.Body.Close()
  body, readErr := io.ReadAll(r.Body)

  if readErr != nil {
    reqObj.Set("body", "")
  } else {
    reqObj.Set("body", string(body))
  }

  headersObj := http.vm.NewObject()
  for key, values := range r.Header {
    if len(values) == 1 {
      headersObj.Set(key, values[0])
    } else if len(values) > 1 {
      headersObj.Set(key, values)
    }
  }
  reqObj.Set("headers", headersObj)

  return reqObj
}

func (http *HTTP) createServer(call goja.FunctionCall) goja.Value {
  if len(call.Arguments) < 1 {
    panic(http.vm.ToValue("createServer requires a request handler function"))
  }

  requestHandler, ok := goja.AssertFunction(call.Arguments[0])
  if !ok {
    panic(http.vm.ToValue("argument must be a function"))
  }

  serverObj := http.vm.NewObject()

  goServer := &netHttp.Server{
    Handler: netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
      reqObj := http.createRequestObject(r)
      resObj := http.vm.NewObject()

      resObj.Set("statusCode", 200)

      resObj.Set("setHeader", func(call goja.FunctionCall) goja.Value {
        if len(call.Arguments) < 2 {
          panic(http.vm.ToValue("setHeader requires a name and value"))
        }
        headerName := call.Arguments[0].String()
        headerValue := call.Arguments[1].String()
        w.Header().Set(headerName, headerValue)

        return goja.Undefined()
      })

      resObj.Set("end", func(call goja.FunctionCall) goja.Value {
        statusCode := 200
        if statusVal := resObj.Get("statusCode"); statusVal != nil && !goja.IsUndefined(statusVal) {
          statusCode = int(statusVal.ToInteger())
        }
        w.WriteHeader(statusCode)

        if len(call.Arguments) > 0 && !goja.IsUndefined(call.Arguments[0]) {
          data := call.Arguments[0].String()
          w.Write([]byte(data))
        }

        return goja.Undefined()
      })

      requestHandler(goja.Undefined(), reqObj, resObj)
    }),
  }

  serverObj.Set("listen", func(call goja.FunctionCall) goja.Value {
    if len(call.Arguments) < 1 {
      panic(http.vm.ToValue("listen requires a port number"))
    }

    port := call.Arguments[0].String()
    addr := ":" + port

    var callback goja.Callable
    if len(call.Arguments) > 1 {
      callback, _ = goja.AssertFunction(call.Arguments[1])
    }

    goServer.Addr = addr

    http.eventLoop.KeepAlive()

    go func() {
      err := goServer.ListenAndServe()
      if err != nil && err != netHttp.ErrServerClosed {
        fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
      }
    }()

    if callback != nil {
      callback(goja.Undefined())
    }
    
    return goja.Undefined()
  })

  serverObj.Set("websocket", func(call goja.FunctionCall) goja.Value {
    if len(call.Arguments) < 2 {
      panic(http.vm.ToValue("websocket requires a url and an object with callback functions"))
    }

    wsPath := call.Arguments[0].String()
    if !strings.HasPrefix(wsPath, "/") {
      panic(http.vm.NewTypeError("websocket path must start with /"))
    }

    callbackObj := call.Arguments[1].ToObject(http.vm)

    if isArray := callbackObj.Get("constructor").String() == "Array"; isArray || callbackObj == nil {
      panic(http.vm.NewTypeError("second argument must be an object"))
    }

    var onOpen, onMessage, onClose, onError goja.Callable

    if openCb := callbackObj.Get("open"); openCb != nil && !goja.IsUndefined(openCb) {
      onOpen, _ = goja.AssertFunction(openCb)
    }

    if messageCb := callbackObj.Get("message"); messageCb != nil && !goja.IsUndefined(messageCb) {
      onMessage, _ = goja.AssertFunction(messageCb)
    }

    if closeCb := callbackObj.Get("close"); closeCb != nil && !goja.IsUndefined(closeCb) {
      onClose, _ = goja.AssertFunction(closeCb)
    }

    if errorCb := callbackObj.Get("error"); errorCb != nil && !goja.IsUndefined(errorCb) {
      onError, _ = goja.AssertFunction(errorCb)
    }

    upgrader := websocket.Upgrader{
      CheckOrigin: func(r *netHttp.Request) bool {
        return true
      },
    }

    mux, ok := goServer.Handler.(*netHttp.ServeMux)
    if !ok {
      mux = netHttp.NewServeMux()
      oldHandler := goServer.Handler
      mux.HandleFunc("/", func(w netHttp.ResponseWriter, r *netHttp.Request) {
        oldHandler.ServeHTTP(w, r)
      })

      goServer.Handler = mux
    }

    mux.HandleFunc(wsPath, func(w netHttp.ResponseWriter, r *netHttp.Request) {
      conn, err := upgrader.Upgrade(w, r, nil)
      if err != nil {
        if onError != nil {
          errMsg := err.Error() // capture data explicitly
          http.eventLoop.ScheduleTask(&event.Task{
            Callback: func() {
              onError(goja.Undefined(), http.vm.ToValue(errMsg))
            },
          })
        }
        return
      }

      wsObj := http.vm.NewObject()
      
      const (
        wsConnecting = 0
        wsOpen       = 1
        wsClosing    = 2
        wsClosed     = 3
      )
      
      var writeMu sync.Mutex
      var state int = wsOpen

      // Helper to update state and sync to property
      // Call this ONLY while holding writeMu lock
      updateState := func(newState int) {
        state = newState
        wsObj.Set("readyState", state)
      }
      
      // Initialize readyState property
      writeMu.Lock()
      wsObj.Set("readyState", state)
      writeMu.Unlock()

      wsObj.Set("CONNECTING", wsConnecting)
      wsObj.Set("OPEN", wsOpen)
      wsObj.Set("CLOSING", wsClosing)
      wsObj.Set("CLOSED", wsClosed)

      wsObj.Set("send", func(call goja.FunctionCall) goja.Value {
        if len(call.Arguments) < 1 {
          panic(http.vm.ToValue("send requires a message"))
        }

        writeMu.Lock()
        currentState := state
        writeMu.Unlock()

        if currentState != wsOpen {
          panic(http.vm.ToValue("websocket connection is not open"))
        }

        message := call.Arguments[0].String()

        writeMu.Lock()
        err := conn.WriteMessage(websocket.TextMessage, []byte(message))
        writeMu.Unlock()

        if err != nil && onError != nil {
          errMsg := err.Error()
          http.eventLoop.ScheduleTask(&event.Task{
            Callback: func() {
              onError(goja.Undefined(), http.vm.ToValue(errMsg))
            },
          })
        }

        return goja.Undefined()
      })

      wsObj.Set("close", func(call goja.FunctionCall) goja.Value {
        writeMu.Lock()
        if state == wsOpen {
          updateState(wsClosing)
        }
        writeMu.Unlock()

        return goja.Undefined()
      })

      if onOpen != nil {
        http.eventLoop.ScheduleTask(&event.Task{
          Callback: func() {
            onOpen(goja.Undefined(), wsObj)
          },
        })
      }

      go func() {
        defer func() {
          writeMu.Lock()
          updateState(wsClosed)
          writeMu.Unlock()
          conn.Close()
        }()

        for {
          messageType, message, err := conn.ReadMessage()

          if err != nil {
            if onError != nil {
              errMsg := err.Error()
              http.eventLoop.ScheduleTask(&event.Task{
                Callback: func() {
                  onError(goja.Undefined(), http.vm.ToValue(errMsg))
                },
              })
            }
            break
          }

          if onMessage != nil {
            var msgData any
            if messageType == websocket.TextMessage {
              msgData = string(message)
            } else {
              msgData = message
            }

            // capture data explicitly in the closure
            capturedData := msgData
            capturedType := messageType

            http.eventLoop.ScheduleTask(&event.Task{
              Callback: func() {
                msgObj := http.vm.NewObject()
                msgObj.Set("data", capturedData)
                msgObj.Set("type", capturedType)
                onMessage(goja.Undefined(), msgObj)
              },
            })
          }
        }

        writeMu.Lock()
        if state != wsClosed {
          updateState(wsClosing)
        }
        writeMu.Unlock()

        if onClose != nil {
          http.eventLoop.ScheduleTask(&event.Task{
            Callback: func() {
              onClose(goja.Undefined())
            },
          })
        }
      }()
    })

    return goja.Undefined()
  })

  return serverObj
}

