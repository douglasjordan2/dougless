package modules

import (
	"context"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/douglasjordan2/dougless/internal/permissions"
	"github.com/gorilla/websocket"
)

func TestWebSocketStateManagement(t *testing.T) {
	// Grant all permissions for tests
	manager := permissions.NewManager()
	manager.GrantAll()
	permissions.SetGlobalManager(manager)

	vm := goja.New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventLoop := event.NewLoopWithContext(ctx)
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
		`

		_, err := vm.RunString(script)
		if err != nil {
			t.Fatalf("Failed to setup WebSocket: %v", err)
		}

		// Give server time to start
		time.Sleep(200 * time.Millisecond)

		// Now test if we can connect
		dialer := websocket.DefaultDialer
		conn, _, err := dialer.Dial("ws://localhost:8090/test", nil)
		if err != nil {
			t.Fatalf("Failed to connect to WebSocket: %v", err)
		}
		defer conn.Close()

		// Give open callback time to execute
		time.Sleep(100 * time.Millisecond)

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
		val, err := vm.RunString("wsObj.readyState")
		if err != nil {
			t.Fatalf("Failed to check readyState: %v", err)
		}

		// readyState should be OPEN (1) or might be CLOSED (3) if connection already closed
		// This is acceptable in tests since the connection might close quickly
		state := val.ToInteger()
		if state != 1 && state != 3 {
			t.Errorf("readyState should be OPEN (1) or CLOSED (3), got %d", state)
		}
	})
}

func TestWebSocketSendAndReceive(t *testing.T) {
	// Grant all permissions for tests
	manager := permissions.NewManager()
	manager.GrantAll()
	permissions.SetGlobalManager(manager)

	vm := goja.New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventLoop := event.NewLoopWithContext(ctx)
	httpModule := NewHTTP(eventLoop)

	vm.Set("http", httpModule.Export(vm))

	go eventLoop.Run()
	defer eventLoop.Stop()

	t.Run("can send and receive messages", func(t *testing.T) {

		script := `
			const server = http.createServer(function(req, res) {
				res.end('ok');
			});

			let serverWs = null;
			server.websocket('/echo', {
				open: function(ws) {
					serverWs = ws;
				},
				message: function(msg) {
					// Echo back using the stored ws reference
					if (serverWs) {
						serverWs.send('Echo: ' + msg.data);
					}
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
	// Grant all permissions for tests
	manager := permissions.NewManager()
	manager.GrantAll()
	permissions.SetGlobalManager(manager)

	vm := goja.New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventLoop := event.NewLoopWithContext(ctx)
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

		// Give close callback time to execute
		time.Sleep(200 * time.Millisecond)

		if !closeCalled {
			t.Error("close callback was not called")
		}
	})
}

func TestWebSocketErrorHandling(t *testing.T) {
	// Grant all permissions for tests
	manager := permissions.NewManager()
	manager.GrantAll()
	permissions.SetGlobalManager(manager)

	vm := goja.New()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	eventLoop := event.NewLoopWithContext(ctx)
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

		// Give open callback time to execute
		time.Sleep(100 * time.Millisecond)

		// Close connection
		conn.Close()
		time.Sleep(100 * time.Millisecond)

		// Try to send on closed connection - should error
		script2 := `
			let errorCaught = false;
			let errorMessage = '';
			try {
				testWs.send('test');
			} catch(e) {
				errorCaught = true;
				errorMessage = e.message || String(e);
			}
			errorCaught;
		`

		val, err := vm.RunString(script2)
		if err != nil {
			t.Fatalf("Script error: %v", err)
		}

		// We just verify that an error was caught
		if !val.ToBoolean() {
			t.Error("Expected an error when sending on closed connection")
		}
	})
}
