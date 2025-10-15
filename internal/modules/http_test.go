package modules

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/dop251/goja"

	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/douglasjordan2/dougless/internal/permissions"
)

// helper to set a fresh permissions manager and restore after test
func withFreshPermissions(t *testing.T) func() {
	t.Helper()
	original := permissions.GetManager()
	permissions.SetGlobalManager(permissions.NewManager())
	return func() { permissions.SetGlobalManager(original) }
}

func TestHTTPGet_Allowed_SendsRequest(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("X-Test", "1")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	// Grant network access to localhost (covers 127.0.0.1 and ::1)
	permissions.GetManager().GrantNet([]string{"localhost"})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	httpMod := NewHTTP(loop)
	httpObj := httpMod.Export(vm)

	getFn, ok := goja.AssertFunction(httpObj.ToObject(vm).Get("get"))
	if !ok {
		t.Fatalf("http.get is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
		gotRes *goja.Object
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		if len(call.Arguments) > 1 && !goja.IsUndefined(call.Arguments[1]) {
			gotRes = call.Arguments[1].ToObject(vm)
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := getFn(goja.Undefined(), vm.ToValue(ts.URL), cb)
	if err != nil {
		t.Fatalf("calling http.get failed: %v", err)
	}

	cbWg.Wait()

	if atomic.LoadInt32(&hits) != 1 {
		t.Fatalf("expected 1 request to be sent, got %d", hits)
	}

	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error on success, got: %v", gotErr)
	}
	if gotRes == nil {
		t.Fatalf("expected response object, got nil")
	}

	statusCode := gotRes.Get("statusCode").ToInteger()
	if statusCode != 200 {
		t.Fatalf("expected statusCode 200, got %d", statusCode)
	}
	body := gotRes.Get("body").String()
	if body != "ok" {
		t.Fatalf("expected body 'ok', got %q", body)
	}
}

func TestHTTPPost_Allowed_SendsRequest(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("created"))
	}))
	defer ts.Close()

	// Grant network access to localhost
	permissions.GetManager().GrantNet([]string{"localhost"})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	httpMod := NewHTTP(loop)
	httpObj := httpMod.Export(vm)

	postFn, ok := goja.AssertFunction(httpObj.ToObject(vm).Get("post"))
	if !ok {
		t.Fatalf("http.post is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
		gotRes *goja.Object
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		if len(call.Arguments) > 1 && !goja.IsUndefined(call.Arguments[1]) {
			gotRes = call.Arguments[1].ToObject(vm)
		}
		cbWg.Done()
		return goja.Undefined()
	})

	payload := map[string]any{"hello": "world"}

	_, err := postFn(goja.Undefined(), vm.ToValue(ts.URL), vm.ToValue(payload), cb)
	if err != nil {
		t.Fatalf("calling http.post failed: %v", err)
	}

	cbWg.Wait()

	if atomic.LoadInt32(&hits) != 1 {
		t.Fatalf("expected 1 request to be sent, got %d", hits)
	}

	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error on success, got: %v", gotErr)
	}
	if gotRes == nil {
		t.Fatalf("expected response object, got nil")
	}

	statusCode := gotRes.Get("statusCode").ToInteger()
	if statusCode != 201 {
		t.Fatalf("expected statusCode 201, got %d", statusCode)
	}
	body := gotRes.Get("body").String()
	if body != "created" {
		t.Fatalf("expected body 'created', got %q", body)
	}
}

func TestHTTPGet_Denied_NoRequestSent(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	httpMod := NewHTTP(loop)
	httpObj := httpMod.Export(vm)

	getFn, ok := goja.AssertFunction(httpObj.ToObject(vm).Get("get"))
	if !ok {
		t.Fatalf("http.get is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
		gotRes goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		if len(call.Arguments) > 1 {
			gotRes = call.Arguments[1]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	// No network permission granted; should return early and not hit server
	_, err := getFn(goja.Undefined(), vm.ToValue(ts.URL), cb)
	if err != nil {
		t.Fatalf("calling http.get failed: %v", err)
	}

	cbWg.Wait()

	if atomic.LoadInt32(&hits) != 0 {
		t.Fatalf("expected no requests to be sent, got %d", hits)
	}

	if goja.IsUndefined(gotErr) || goja.IsNull(gotErr) {
		t.Fatalf("expected an error value when permission denied")
	}
	// data should be undefined on denied
	if !goja.IsUndefined(gotRes) {
		t.Fatalf("expected data to be undefined on denied, got: %v", gotRes)
	}
}

func TestHTTPCreateServer_ListenDenied_PanicsWithoutPermission(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	httpMod := NewHTTP(loop)
	httpObj := httpMod.Export(vm)

	createServerFn, ok := goja.AssertFunction(httpObj.ToObject(vm).Get("createServer"))
	if !ok {
		t.Fatalf("http.createServer is not a function")
	}

	// handler just ends the response
	handler := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		// noop
		return goja.Undefined()
	})

	serverVal, err := createServerFn(goja.Undefined(), handler)
	if err != nil {
		t.Fatalf("createServer call failed: %v", err)
	}

	listenFn, ok := goja.AssertFunction(serverVal.ToObject(vm).Get("listen"))
	if !ok {
		t.Fatalf("server.listen is not a function")
	}

	// No PermissionNet granted; should return an exception error
	_, callErr := listenFn(goja.Undefined(), vm.ToValue("0"))
	if callErr == nil {
		t.Fatalf("expected listen to fail with an error when PermissionNet is denied")
	}
}

func TestHTTPCreateServer_Allowed_ListensAndServes(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	// Allow localhost
	permissions.GetManager().GrantNet([]string{"localhost"})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	httpMod := NewHTTP(loop)
	httpObj := httpMod.Export(vm)

	createServerFn, ok := goja.AssertFunction(httpObj.ToObject(vm).Get("createServer"))
	if !ok {
		t.Fatalf("http.createServer is not a function")
	}

	// handler: echo simple body
	handler := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		res := call.Arguments[1].ToObject(vm)
		endFn, _ := goja.AssertFunction(res.Get("end"))
		_, _ = endFn(goja.Undefined(), vm.ToValue("pong"))
		return goja.Undefined()
	})

	serverVal, err := createServerFn(goja.Undefined(), handler)
	if err != nil {
		t.Fatalf("createServer call failed: %v", err)
	}

	listenFn, ok := goja.AssertFunction(serverVal.ToObject(vm).Get("listen"))
	if !ok {
		t.Fatalf("server.listen is not a function")
	}

	// listen on port 0 and bind to loopback; wait for callback
	cbDone := make(chan struct{})
	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		close(cbDone)
		return goja.Undefined()
	})

	_, callErr := listenFn(goja.Undefined(), vm.ToValue("0"), vm.ToValue("127.0.0.1"), cb)
	if callErr != nil {
		t.Fatalf("listen failed: %v", callErr)
	}

	<-cbDone

	addrVal := serverVal.ToObject(vm).Get("address")
	if goja.IsUndefined(addrVal) || goja.IsNull(addrVal) {
		t.Fatalf("expected server address to be set")
	}
	addr := addrVal.String()

	// Do a request
	resp, err := http.Get("http://" + addr)
	if err != nil {
		t.Fatalf("http GET failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	// Close the server
	closeFn, ok := goja.AssertFunction(serverVal.ToObject(vm).Get("close"))
	if !ok {
		t.Fatalf("server.close is not a function")
	}
	_, _ = closeFn(goja.Undefined())
}

func TestHTTPPost_Denied_NoRequestSent(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	var hits int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	defer ts.Close()

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	httpMod := NewHTTP(loop)
	httpObj := httpMod.Export(vm)

	postFn, ok := goja.AssertFunction(httpObj.ToObject(vm).Get("post"))
	if !ok {
		t.Fatalf("http.post is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
		gotRes goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		if len(call.Arguments) > 1 {
			gotRes = call.Arguments[1]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	payload := map[string]any{"hello": "world"}

	_, err := postFn(goja.Undefined(), vm.ToValue(ts.URL), vm.ToValue(payload), cb)
	if err != nil {
		t.Fatalf("calling http.post failed: %v", err)
	}

	cbWg.Wait()

	if atomic.LoadInt32(&hits) != 0 {
		t.Fatalf("expected no requests to be sent, got %d", hits)
	}

	if goja.IsUndefined(gotErr) || goja.IsNull(gotErr) {
		t.Fatalf("expected an error value when permission denied")
	}
	if !goja.IsUndefined(gotRes) {
		t.Fatalf("expected data to be undefined on denied, got: %v", gotRes)
	}
}
