package modules

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	netHttp "net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dop251/goja"
	"github.com/gorilla/websocket"

	"github.com/douglasjordan2/dougless/internal/permissions"
	"github.com/douglasjordan2/dougless/internal/future"
)

type HTTP struct {
	vm        *goja.Runtime 
  taskQueue chan func()
  runtime   RuntimeKeepAlive
}

func (http *HTTP) SetRuntime(rt RuntimeKeepAlive) {
  http.runtime = rt
}

func NewHTTP(vm *goja.Runtime) *HTTP {
  h := &HTTP{
    vm:        vm,
    taskQueue: make(chan func(), 100),
  }

  go h.runTaskQueue() // dedicated goroutine for VM tasks

  return h
}

func (http *HTTP) runTaskQueue() {
  for task := range http.taskQueue {
    task()
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

func (http *HTTP) extractHost(urlStr string) string {
	urlStr = strings.TrimPrefix(urlStr, "http://")
	urlStr = strings.TrimPrefix(urlStr, "https://")

	parts := strings.SplitN(urlStr, "/", 2)
	return parts[0]
}

func (http *HTTP) argCheck(call goja.FunctionCall, size int, errMsg string) {
	if len(call.Arguments) < size {
		panic(http.vm.ToValue(errMsg))
	}
}

func (http *HTTP) hasNetPermissions(url string, ctx context.Context) (string, bool) {
  host := http.extractHost(url)
  mgr := permissions.GetManager()
  canNet := permissions.PermissionNet

  canAccess := true
  if !mgr.CheckWithPrompt(ctx, canNet, host) {
    canAccess = false
  }

  return host, canAccess
}

func (http *HTTP) getHeaders(resp *netHttp.Response) map[string]any {
  headers := make(map[string]any)
  for key, values := range resp.Header {
    if len(values) == 1 {
      headers[key] = values[0]
    } else if len(values) > 1 {
      headers[key] = values
    }
  }
  return headers
}


func createProxy(vm *goja.Runtime, future *future.Future) goja.Value {
  obj := vm.NewObject()

  getter := func(key string) any {
    result, ok := future.MustGet().(map[string]any)
    if !ok {
      return nil
    }
    return result[key]
  }

  obj.DefineAccessorProperty("status",
    vm.ToValue(func() any { return getter("status") }), // getter
    nil, // setter
    goja.FLAG_FALSE, // is writeable?
    goja.FLAG_TRUE) // is enumerable?

  obj.DefineAccessorProperty("body",
    vm.ToValue(func() any { return getter("body") }), nil, goja.FLAG_FALSE, goja.FLAG_TRUE)

  obj.DefineAccessorProperty("headers",
    vm.ToValue(func() any { return getter("headers") }), nil, goja.FLAG_FALSE, goja.FLAG_TRUE)

  return obj
}


func (http *HTTP) get(call goja.FunctionCall) goja.Value {
  http.argCheck(call, 1, "GET requires a URL")

	url := call.Arguments[0].String()

	f := future.NewFuture(func() (any, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		host, canAccess := http.hasNetPermissions(url, ctx)
		if !canAccess {
			return nil, fmt.Errorf("permission denied for %s", host)
		}

		resp, err := netHttp.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, readErr
		}

		headers := http.getHeaders(resp)

		return map[string]any{
			"statusCode":     resp.StatusCode,
			"statusText": resp.Status,
			"body":       string(body),
			"headers":    headers,
		}, nil
	})

	return createProxy(http.vm, f)
}

func (http *HTTP) post(call goja.FunctionCall) goja.Value {
  http.argCheck(call, 2, "POST requires a URL and a payload")

	url := call.Arguments[0].String()
	payload := call.Arguments[1].Export()

	contentType := "application/json"
	dataMap, isMap := payload.(map[string]any)

	if isMap {
		if ct, exists := dataMap["contentType"]; exists {
			contentType = ct.(string)
			delete(dataMap, "contentType")
			payload = dataMap
		}
	}

  f := future.NewFuture(func() (any, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    host, canAccess := http.hasNetPermissions(url, ctx) 
    if !canAccess {
      return nil, fmt.Errorf("permission denied for %s", host)
    }

    jsonBytes, marshalErr := json.Marshal(payload)
    if marshalErr != nil {
      return nil, marshalErr
    }

    body := bytes.NewBuffer(jsonBytes)
    resp, err := netHttp.Post(url, contentType, body)
    if err != nil {
      return nil, err
    }
    defer resp.Body.Close()

    respBody, readErr := io.ReadAll(resp.Body)
    if readErr != nil {
      return nil, readErr
    }

		headers := http.getHeaders(resp)

    return map[string]any{
      "statusCode":     resp.StatusCode,
      "statusText": resp.Status,
      "body":       string(respBody),
      "headers":    headers,
    }, nil
  })

	return createProxy(http.vm, f)
}

func (http *HTTP) createRequestObject(r *netHttp.Request) goja.Value {
	reqObj := http.vm.NewObject()

	reqObj.Set("method", r.Method)
	reqObj.Set("url", r.URL.String())

	body, readErr := io.ReadAll(r.Body)
	r.Body.Close()

	if readErr != nil {
		reqObj.Set("body", "")
	} else {
		reqObj.Set("body", string(body))
    r.Body = io.NopCloser(bytes.NewBuffer(body))
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

  type responseState struct {
    statusCode int
    headers    map[string]string
    body       string
    mu         sync.Mutex
  }

	goServer := &netHttp.Server{
		Handler: netHttp.HandlerFunc(func(w netHttp.ResponseWriter, r *netHttp.Request) {
      done := make(chan struct{})
      state := &responseState{
        statusCode: 200,
        headers:    make(map[string]string),
      }

      http.taskQueue <- func() {
        defer close(done)

        reqObj := http.createRequestObject(r)
        resObj := http.vm.NewObject()

        resObj.Set("statusCode", 200)

        resObj.Set("setHeader", func(call goja.FunctionCall) goja.Value {
          if len(call.Arguments) < 2 {
            panic(http.vm.ToValue("setHeader requires a name and value"))
          }
          headerName := call.Arguments[0].String()
          headerValue := call.Arguments[1].String()

          state.mu.Lock()
          state.headers[headerName] = headerValue
          state.mu.Unlock()

          return goja.Undefined()
        })

        resObj.Set("writeHead", func(call goja.FunctionCall) goja.Value {
          if len(call.Arguments) < 1 {
            panic(http.vm.ToValue("writeHead requires a status code"))
          }
          statusCode := int(call.Arguments[0].ToInteger())
          
          state.mu.Lock()
          state.statusCode = statusCode

          if len(call.Arguments) > 1 && !goja.IsUndefined(call.Arguments[1]) {
            headersObj := call.Arguments[1].ToObject(http.vm)
            for _, key := range headersObj.Keys() {
              val := headersObj.Get(key) 
              if !goja.IsUndefined(val) {
                state.headers[key] = val.String()
              }
            }
          }
          state.mu.Unlock()

          return goja.Undefined()
        })

        resObj.Set("end", func(call goja.FunctionCall) goja.Value {
          state.mu.Lock()
          if statusVal := resObj.Get("statusCode"); statusVal != nil && !goja.IsUndefined(statusVal) {
            state.statusCode = int(statusVal.ToInteger())
          }

          if len(call.Arguments) > 0 && !goja.IsUndefined(call.Arguments[0]) {
            state.body = call.Arguments[0].String()
          }
          state.mu.Unlock()

          return goja.Undefined()
        })

        requestHandler(goja.Undefined(), reqObj, resObj)
      }
      
      select {
      case <-done:
        state.mu.Lock()
        for name, value := range state.headers {
          w.Header().Set(name, value)
        }
        w.WriteHeader(state.statusCode)
        if state.body != "" {
          w.Write([]byte(state.body))
        }
        state.mu.Unlock()
      case <-time.After(30 * time.Second):
        w.WriteHeader(netHttp.StatusGatewayTimeout)
        w.Write([]byte("Request handler timeout"))
      }
		}),
	}

	serverObj.Set("listen", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) < 1 {
			panic(http.vm.ToValue("listen requires a port number"))
		}

		port := call.Arguments[0].String()
		bindAddr := "0.0.0.0"
		argOffset := 1

		if len(call.Arguments) > 1 {
			if _, ok := goja.AssertFunction(call.Arguments[1]); !ok {
				bindAddr = call.Arguments[1].String()
				argOffset = 2
			}
		}

		var callback goja.Callable
		if len(call.Arguments) > argOffset {
			callback, _ = goja.AssertFunction(call.Arguments[argOffset])
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		permHost := bindAddr + ":" + port
		mgr := permissions.GetManager()
		canNet := permissions.PermissionNet
		if !mgr.CheckWithPrompt(ctx, canNet, permHost) {
			errMsg := mgr.ErrorMessage(canNet, permHost)
			panic(http.vm.ToValue(errMsg))
		}

		ln, err := net.Listen("tcp", bindAddr+":"+port)
		if err != nil {
			panic(http.vm.ToValue(err.Error()))
		}

		goServer.Addr = ln.Addr().String()
		serverObj.Set("address", goServer.Addr)

    done := http.runtime.KeepAlive()
		go func() {
      defer done()
			err := goServer.Serve(ln)
			if err != nil && err != netHttp.ErrServerClosed {
				fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			}
		}()

		if callback != nil {
			callback(goja.Undefined())
		}

		return goja.Undefined()
	})

	serverObj.Set("close", func(call goja.FunctionCall) goja.Value {
		_ = goServer.Close()
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
					errMsg := err.Error() 
					http.taskQueue <- func() {
            onError(goja.Undefined(), http.vm.ToValue(errMsg))
					}
				}
				return
			}

			const (
				wsConnecting = 0
				wsOpen       = 1
				wsClosing    = 2
				wsClosed     = 3
			)

			var writeMu sync.Mutex
			var state int = wsOpen
			ctx, cancel := context.WithCancel(context.Background())

			// Create and setup WebSocket object in VM-safe goroutine
			wsObjChan := make(chan *goja.Object)
			http.taskQueue <- func() {
				wsObj := http.vm.NewObject()
				
				wsObj.Set("readyState", wsOpen)
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
						http.taskQueue <- func() {
							onError(goja.Undefined(), http.vm.ToValue(errMsg))
						}
					}

					return goja.Undefined()
				})

				wsObj.Set("close", func(call goja.FunctionCall) goja.Value {
					writeMu.Lock()
					if state == wsOpen || state == wsConnecting {
						state = wsClosing
						closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
						conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
						cancel()
					}
					writeMu.Unlock()

					http.taskQueue <- func() {
						writeMu.Lock()
						currentState := state
						writeMu.Unlock()
						wsObj.Set("readyState", currentState)
					}

					return goja.Undefined()
				})
				
				wsObjChan <- wsObj
			}
			wsObj := <-wsObjChan

			if onOpen != nil {
				http.taskQueue <- func() {
          onOpen(goja.Undefined(), wsObj)
				}
			}

      done := http.runtime.KeepAlive()
			go func() {
				defer func() {
          done()
					cancel() 

					writeMu.Lock()
					state = wsClosed
					writeMu.Unlock()

          http.taskQueue <- func() {
            wsObj.Set("readyState", wsClosed)
          }

					conn.Close()
				}()

				readDone := make(chan struct{})
        done := http.runtime.KeepAlive()
				go func() {
          defer done()
					<-ctx.Done()
					conn.SetReadDeadline(time.Now()) 
					close(readDone)
				}()

				for {
					select {
					case <-ctx.Done():
						return
					default:
					}

					messageType, message, err := conn.ReadMessage()

					if err != nil {
						select {
						case <-ctx.Done():
						default:
							if onError != nil {
								errMsg := err.Error()
								http.taskQueue <- func() {
                  onError(goja.Undefined(), http.vm.ToValue(errMsg))
								}
							}
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

						capturedData := msgData
						capturedType := messageType

						http.taskQueue <- func() {
              msgObj := http.vm.NewObject()
              msgObj.Set("data", capturedData)
              msgObj.Set("type", capturedType)
              onMessage(goja.Undefined(), msgObj)
						}
					}
				}

				var shouldUpdateState bool
				writeMu.Lock()
				if state != wsClosed {
					state = wsClosing
					shouldUpdateState = true
				}
				writeMu.Unlock()

				if shouldUpdateState {
          http.taskQueue <- func() {
            wsObj.Set("readyState", wsClosing)
          }
				}

				if onClose != nil {
					http.taskQueue <- func() {
            onClose(goja.Undefined())
					}
				}
			}()
		})

		return goja.Undefined()
	})

	return serverObj
}
