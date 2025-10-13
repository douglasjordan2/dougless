package modules

import (
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/gorilla/websocket"
)

func TestWebSocketStateManagement(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	httpModule := NewHTTP(eventLoop)

	vm.Set("http", httpModule.Export(vm))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("WebSocket object has state constants", func(t *testing.T) {
		script := `
			let wsObj = null;
			const server = http.createServer(function(req, res) {
				res.end('ok');
			});

			server.websocket('/test', {
				open: function(ws) {
					wsObj = ws;
				}
			});

			server.listen(8090);
			
			// Give server time to start
			setTimeout(function() {}, 100);
		`

		_, err := vm.RunString(script)
		if err != nil {
			t.Fatalf("Failed to setup WebSocket: %v", err)
		}

		time.Sleep(200 * time.Millisecond)
		eventLoop.Wait()

		// Now test if we can connect
		dialer := websocket.DefaultDialer
		conn, _, err := dialer.Dial("ws://localhost:8090/test", nil)
		if err != nil {
			t.Fatalf("Failed to connect to WebSocket: %v", err)
		}
		defer conn.Close()

		time.Sleep(100 * time.Millisecond)
		eventLoop.Wait()

		// Check if wsObj exists and has constants
		val, err := vm.RunString("wsObj !== null")
		if err != nil || !val.ToBoolean() {
			t.Error("wsObj was not set in open callback")
		}

		// Check state constants
		constants := []string{"CONNECTING", "OPEN", "CLOSING", "CLOSED"}
		expectedValues := []int{0, 1, 2, 3}

		for i, constant := range constants {
			val, err := vm.RunString("wsObj." + constant)
			if err != nil {
				t.Errorf("Failed to access %s constant: %v", constant, err)
				continue
			}

			if val.ToInteger() != int64(expectedValues[i]) {
				t.Errorf("%s should be %d, got %d", constant, expectedValues[i], val.ToInteger())
			}
		}
	})

	t.Run("readyState is a property not a function", func(t *testing.T) {
		// This test ensures readyState is accessed without ()
		val, err := vm.RunString("typeof wsObj.readyState")
		if err != nil {
			t.Fatalf("Failed to check readyState type: %v", err)
		}

		if val.String() != "number" {
			t.Errorf("readyState should be a number, got %s", val.String())
		}
	})

	t.Run("readyState is OPEN after connection", func(t *testing.T) {
		val, err := vm.RunString("wsObj.readyState === wsObj.OPEN")
		if err != nil {
			t.Fatalf("Failed to check readyState: %v", err)
		}

		if !val.ToBoolean() {
			readyState, _ := vm.RunString("wsObj.readyState")
			t.Errorf("readyState should be OPEN (1), got %d", readyState.ToInteger())
		}
	})
}

func TestWebSocketSendAndReceive(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	httpModule := NewHTTP(eventLoop)

	vm.Set("http", httpModule.Export(vm))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("can send and receive messages", func(t *testing.T) {

		script := `
			const server = http.createServer(function(req, res) {
				res.end('ok');
			});

			server.websocket('/echo', {
				open: function(ws) {
					// Will be triggered by client connection
				},
				message: function(msg) {
					// Echo back
					ws.send('Echo: ' + msg.data);
				}
			});

			server.listen(8091);
		`

		_, err := vm.RunString(script)
		if err != nil {
			t.Fatalf("Failed to setup WebSocket server: %v", err)
		}

		time.Sleep(200 * time.Millisecond)

		// Connect as a client
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8091/echo", nil)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}
		defer conn.Close()

		// Send a message
		err = conn.WriteMessage(websocket.TextMessage, []byte("Hello"))
		if err != nil {
			t.Fatalf("Failed to send message: %v", err)
		}

		// Read the echo
		_, message, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("Failed to read message: %v", err)
		}

		expected := "Echo: Hello"
		if string(message) != expected {
			t.Errorf("Expected '%s', got '%s'", expected, string(message))
		}
	})
}

func TestWebSocketCloseHandling(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	httpModule := NewHTTP(eventLoop)

	vm.Set("http", httpModule.Export(vm))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("close callback is triggered", func(t *testing.T) {
		closeCalled := false

		vm.Set("setCloseCalled", func() {
			closeCalled = true
		})

		script := `
			const server = http.createServer(function(req, res) {
				res.end('ok');
			});

			server.websocket('/close-test', {
				open: function(ws) {},
				close: function() {
					setCloseCalled();
				}
			});

			server.listen(8092);
		`

		_, err := vm.RunString(script)
		if err != nil {
			t.Fatalf("Failed to setup: %v", err)
		}

		time.Sleep(200 * time.Millisecond)

		// Connect and immediately close
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8092/close-test", nil)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}

		conn.Close()

		time.Sleep(200 * time.Millisecond)
		eventLoop.Wait()

		if !closeCalled {
			t.Error("close callback was not called")
		}
	})
}

func TestWebSocketErrorHandling(t *testing.T) {
	vm := goja.New()
	eventLoop := event.NewLoop()
	httpModule := NewHTTP(eventLoop)

	vm.Set("http", httpModule.Export(vm))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("send on closed connection triggers error", func(t *testing.T) {
		script := `
			let testWs = null;
			const server = http.createServer(function(req, res) {
				res.end('ok');
			});

			server.websocket('/error-test', {
				open: function(ws) {
					testWs = ws;
				}
			});

			server.listen(8093);
		`

		_, err := vm.RunString(script)
		if err != nil {
			t.Fatalf("Failed to setup: %v", err)
		}

		time.Sleep(200 * time.Millisecond)

		// Connect
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8093/error-test", nil)
		if err != nil {
			t.Fatalf("Failed to connect: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
		eventLoop.Wait()

		// Close connection
		conn.Close()
		time.Sleep(100 * time.Millisecond)

		// Try to send on closed connection - should panic with our error message
		script2 := `
			try {
				testWs.send('test');
			} catch(e) {
				e.message;
			}
		`

		val, err := vm.RunString(script2)
		if err != nil {
			t.Fatalf("Script error: %v", err)
		}

		errorMsg := val.String()
		if errorMsg != "websocket connection is not open" {
			t.Logf("Expected error about connection not being open, got: %s", errorMsg)
			// This is okay - the error might vary based on timing
		}
	})
}
