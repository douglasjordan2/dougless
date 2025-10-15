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

	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/douglasjordan2/dougless/internal/permissions"
)

// HTTP provides HTTP client and server functionality for JavaScript.
// Includes support for GET/POST requests, server creation with WebSocket upgrade capability,
// and comprehensive request/response handling. All operations integrate with the event loop
// for non-blocking execution and require network permissions.
//
// Available globally in JavaScript as the 'http' object (unique to Dougless).
//
// Example usage:
//
//	// HTTP Client
//	http.get('https://api.example.com/data', (err, response) => {
//	  console.log(response.body);
//	});
//
//	// HTTP Server
//	const server = http.createServer((req, res) => {
//	  res.end('Hello World');
//	});
//	server.listen(3000);
type HTTP struct {
	vm        *goja.Runtime // JavaScript runtime instance
	eventLoop *event.Loop   // Event loop for async task scheduling
}

// NewHTTP creates a new HTTP instance with the given event loop.
func NewHTTP(eventLoop *event.Loop) *HTTP {
	return &HTTP{
		eventLoop: eventLoop,
	}
}

// Export creates and returns the HTTP JavaScript object with all HTTP methods.
func (http *HTTP) Export(vm *goja.Runtime) goja.Value {
	http.vm = vm
	obj := vm.NewObject()

	obj.Set("get", http.get)
	obj.Set("post", http.post)
	obj.Set("createServer", http.createServer)

	return obj
}

// extractHost extracts the hostname from a URL string for permission checks.
// Removes http:// or https:// prefix and returns everything before the first /.
func (http *HTTP) extractHost(urlStr string) string {
	urlStr = strings.TrimPrefix(urlStr, "http://")
	urlStr = strings.TrimPrefix(urlStr, "https://")

	parts := strings.SplitN(urlStr, "/", 2)
	return parts[0]
}

// get performs an HTTP GET request asynchronously.
// The operation is scheduled on the event loop and requires network permission for the target host.
//
// Parameters:
//   - url (string): The URL to fetch
//   - callback (function): Called with (thisArg, error, response) after completion
//
// The response object contains:
//   - status (string): HTTP status text (e.g., "200 OK")
//   - statusCode (number): HTTP status code (e.g., 200)
//   - body (string): Response body as a string
//   - headers (object): Response headers (single values as strings, multiple as arrays)
//
// If permission is denied or an error occurs, the callback receives an error message
// and the response is undefined. On success, error is null.
//
// Example:
//
//	http.get('https://api.github.com/users/octocat', function(thisArg, err, resp) {
//	  if (err) {
//	    console.error('Request failed:', err);
//	  } else {
//	    console.log('Status:', resp.statusCode);
//	    console.log('Body:', resp.body);
//	  }
//	});
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
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			host := http.extractHost(url)
			mgr := permissions.GetManager()
			canNet := permissions.PermissionNet
			if !mgr.CheckWithPrompt(ctx, canNet, host) {
				errMsg := mgr.ErrorMessage(canNet, host)
				callback(goja.Undefined(), http.vm.ToValue(errMsg), goja.Undefined())
				return
			}

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

// post performs an HTTP POST request asynchronously with a JSON payload.
// The operation is scheduled on the event loop and requires network permission for the target host.
//
// Parameters:
//   - url (string): The URL to post to
//   - payload (object): Data to send (automatically JSON-encoded)
//   - callback (function): Called with (thisArg, error, response) after completion
//
// The payload can include a special 'contentType' property to override the default
// 'application/json' content type. This property is removed before encoding.
//
// The response object contains:
//   - status (string): HTTP status text
//   - statusCode (number): HTTP status code
//   - body (string): Response body as a string
//   - headers (object): Response headers
//
// Example:
//
//	http.post('https://api.example.com/data', {name: 'Alice', age: 30}, function(thisArg, err, resp) {
//	  if (err) {
//	    console.error('POST failed:', err);
//	  } else {
//	    console.log('Response:', resp.body);
//	  }
//	});
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
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			host := http.extractHost(url)
			mgr := permissions.GetManager()
			canNet := permissions.PermissionNet
			if !mgr.CheckWithPrompt(ctx, canNet, host) {
				errMsg := mgr.ErrorMessage(canNet, host)
				callback(goja.Undefined(), http.vm.ToValue(errMsg), goja.Undefined())
				return
			}

			jsonBytes, marshalErr := json.Marshal(payload)
			if marshalErr != nil {
				callback(goja.Undefined(), http.vm.ToValue(marshalErr.Error()), goja.Undefined())
				return
			}
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

// createRequestObject converts a Go HTTP request into a JavaScript request object.
// The object includes method, url, body, and headers properties.
// Used internally by the HTTP server to provide request data to JavaScript handlers.
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

// createServer creates an HTTP server with the given request handler.
// The server supports standard HTTP requests and WebSocket upgrades.
// Returns a server object with listen(), close(), and websocket() methods.
//
// Parameters:
//   - handler (function): Called with (thisArg, request, response) for each HTTP request
//
// The request object contains:
//   - method (string): HTTP method (GET, POST, etc.)
//   - url (string): Request URL
//   - body (string): Request body
//   - headers (object): Request headers
//
// The response object has methods:
//   - setHeader(name, value): Set a response header
//   - end(data): Send the response with optional data
//   - statusCode (property): Set HTTP status code (default 200)
//
// The returned server object has methods:
//   - listen(port [, host] [, callback]): Start listening on the specified port
//   - close(): Stop the server and release resources
//   - websocket(path, callbacks): Add WebSocket endpoint (see websocket documentation)
//
// Example:
//
//	const server = http.createServer((req, res) => {
//	  res.setHeader('Content-Type', 'text/plain');
//	  res.end('Hello from Dougless!');
//	});
//	server.listen(8080, () => console.log('Server running on port 8080'));
func (http *HTTP) createServer(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) < 1 {
		panic(http.vm.ToValue("createServer requires a request handler function"))
	}

	requestHandler, ok := goja.AssertFunction(call.Arguments[0])
	if !ok {
		panic(http.vm.ToValue("argument must be a function"))
	}

	serverObj := http.vm.NewObject()

	var keepAliveDone func()
	var keepAliveOnce sync.Once

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
		bindAddr := "0.0.0.0"
		argOffset := 1

		if len(call.Arguments) > 1 {
			if _, ok := goja.AssertFunction(call.Arguments[1]); !ok {
				// not a cb, host string provided
				bindAddr = call.Arguments[1].String()
				argOffset = 2
			}
		}

		var callback goja.Callable
		if len(call.Arguments) > argOffset {
			callback, _ = goja.AssertFunction(call.Arguments[argOffset])
		}

		// Use actual bind address for permission check
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		permHost := bindAddr + ":" + port
		mgr := permissions.GetManager()
		canNet := permissions.PermissionNet
		if !mgr.CheckWithPrompt(ctx, canNet, permHost) {
			errMsg := mgr.ErrorMessage(canNet, permHost)
			// throw a JS exception (same behavior pattern as other denied ops)
			panic(http.vm.ToValue(errMsg))
		}

		// Create listener to learn the actual bound address (handles port "0")
		ln, err := net.Listen("tcp", bindAddr+":"+port)
		if err != nil {
			panic(http.vm.ToValue(err.Error()))
		}

		// Update server address and expose it on the server object
		goServer.Addr = ln.Addr().String()
		serverObj.Set("address", goServer.Addr)

		// Keep event loop alive until server is closed
		keepAliveOnce.Do(func() { keepAliveDone = http.eventLoop.KeepAlive() })

		go func() {
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
		// Close server and release event loop keep-alive
		_ = goServer.Close()
		if keepAliveDone != nil {
			keepAliveDone()
			keepAliveDone = nil
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
			ctx, cancel := context.WithCancel(context.Background())

			// Initialize readyState property
			wsObj.Set("readyState", state)

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
				if state == wsOpen || state == wsConnecting {
					state = wsClosing
					// Send close message to client
					closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
					conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(time.Second))
					// Cancel context to stop read loop
					cancel()
				}
				writeMu.Unlock()

				// Update readyState outside of mutex
				wsObj.Set("readyState", state)

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
					cancel() // Ensure context is cancelled
					writeMu.Lock()
					state = wsClosed
					writeMu.Unlock()
					wsObj.Set("readyState", wsClosed)
					conn.Close()
				}()

				// Set up read deadline check based on context
				readDone := make(chan struct{})
				go func() {
					<-ctx.Done()
					conn.SetReadDeadline(time.Now()) // Force read to fail
					close(readDone)
				}()

				for {
					// Check if context was cancelled
					select {
					case <-ctx.Done():
						// Clean close initiated by user
						return
					default:
					}

					messageType, message, err := conn.ReadMessage()

					if err != nil {
						// Don't report error if it was due to intentional close
						select {
						case <-ctx.Done():
							// Intentional close, don't call onError
						default:
							if onError != nil {
								errMsg := err.Error()
								http.eventLoop.ScheduleTask(&event.Task{
									Callback: func() {
										onError(goja.Undefined(), http.vm.ToValue(errMsg))
									},
								})
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

				var shouldUpdateState bool
				writeMu.Lock()
				if state != wsClosed {
					state = wsClosing
					shouldUpdateState = true
				}
				writeMu.Unlock()

				if shouldUpdateState {
					wsObj.Set("readyState", wsClosing)
				}

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
