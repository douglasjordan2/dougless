package modules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dop251/goja"
	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/douglasjordan2/dougless/internal/permissions"
)

func setupFileTest(t *testing.T) (*goja.Runtime, *FileSystem) {
	// Grant all permissions for tests
	manager := permissions.NewManager()
	manager.GrantAll()
	permissions.SetGlobalManager(manager)

	vm := goja.New()
	loop := event.NewLoop()
	fs := NewFileSystem(loop)

	go loop.Run()

	t.Cleanup(func() {
		loop.Stop()
	})

	return vm, fs
}

func createTestFile(t *testing.T, content string) string {
	tmpfile, err := os.CreateTemp("", "dougless-test-*.txt")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.Remove(tmpfile.Name())
	})

	return tmpfile.Name()
}

func createTestDir(t *testing.T) string {
	tmpdir, err := os.MkdirTemp("", "dougless-test-dir-*")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpdir)
	})

	return tmpdir
}

func TestFileRead(t *testing.T) {
	vm, fs := setupFileTest(t)

	testContent := "Hello from Dougless!"
	testFile := createTestFile(t, testContent)

	fileObj := fs.Export(vm)

	callbackCalled := false
	var readContent string
	var readError string

	callback := func(call goja.FunctionCall) goja.Value {
		callbackCalled = true

		if !goja.IsNull(call.Argument(0)) && !goja.IsUndefined(call.Argument(0)) {
			readError = call.Argument(0).String()
		}

		if !goja.IsUndefined(call.Argument(1)) {
			readContent = call.Argument(1).String()
		}

		return goja.Undefined()
	}

	readFunc, _ := goja.AssertFunction(fileObj.ToObject(vm).Get("read"))
	readFunc(goja.Undefined(), vm.ToValue(testFile), vm.ToValue(callback))

	fs.eventLoop.Wait()

	if !callbackCalled {
		t.Fatal("Callback was not called")
	}

	if readError != "" {
		t.Fatalf("Expected no error, got %s", readError)
	}

	if readContent != testContent {
		t.Fatalf("Expected content %q, got %q!", testContent, readContent)
	}
}

func TestFileWrite(t *testing.T) {
	vm, fs := setupFileTest(t)
	fileObj := fs.Export(vm)

	testDir := createTestDir(t)
	testFile := filepath.Join(testDir, "write-test.txt")
	testContent := "Written by Dougless"

	callbackCalled := false
	var writeError string

	callback := func(call goja.FunctionCall) goja.Value {
		callbackCalled = true
		if !goja.IsNull(call.Argument(0)) && !goja.IsUndefined(call.Argument(0)) {
			writeError = call.Argument(0).String()
		}
		return goja.Undefined()
	}

	writeFunc, _ := goja.AssertFunction(fileObj.ToObject(vm).Get("write"))
	writeFunc(goja.Undefined(), vm.ToValue(testFile), vm.ToValue(testContent), vm.ToValue(callback))

	fs.eventLoop.Wait()

	if !callbackCalled {
		t.Fatal("Callback was not called")
	}

	if writeError != "" {
		t.Fatalf("Expected no error, got: %s", writeError)
	}

	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(content) != testContent {
		t.Fatalf("Expected %q, got %q", testContent, string(content))
	}
}

func TestFileExists(t *testing.T) {
	vm, fs := setupFileTest(t)
	fileObj := fs.Export(vm)

	existingFile := createTestFile(t, "test")

	callbackCalled := false
	var exists bool

	callback := func(call goja.FunctionCall) goja.Value {
		callbackCalled = true
		exists = call.Argument(0).ToBoolean()
		return goja.Undefined()
	}

	existsFunc, _ := goja.AssertFunction(fileObj.ToObject(vm).Get("exists"))
	existsFunc(goja.Undefined(), vm.ToValue(existingFile), vm.ToValue(callback))

	fs.eventLoop.Wait()

	if !callbackCalled {
		t.Fatal("Callback was not called")
	}

	if !exists {
		t.Fatal("Expected file to exist")
	}

	callbackCalled = false
	existsFunc(goja.Undefined(), vm.ToValue("/nonexistent/file.txt"), vm.ToValue(callback))

	fs.eventLoop.Wait()

	if !callbackCalled {
		t.Fatal("Callback was not called for non-existing file")
	}

	if exists {
		t.Fatal("Expected file to not exist")
	}
}

func TestFileMkdir(t *testing.T) {
	vm, fs := setupFileTest(t)
	fileObj := fs.Export(vm)

	testDir := createTestDir(t)
	newDir := filepath.Join(testDir, "new-directory")

	callbackCalled := false
	var mkdirError string

	callback := func(call goja.FunctionCall) goja.Value {
		callbackCalled = true
		if !goja.IsNull(call.Argument(0)) && !goja.IsUndefined(call.Argument(0)) {
			mkdirError = call.Argument(0).String()
		}
		return goja.Undefined()
	}

	mkdirFunc, _ := goja.AssertFunction(fileObj.ToObject(vm).Get("mkdir"))
	mkdirFunc(goja.Undefined(), vm.ToValue(newDir), vm.ToValue(callback))

	fs.eventLoop.Wait()

	if !callbackCalled {
		t.Fatal("Callback was not called")
	}

	if mkdirError != "" {
		t.Fatalf("Expected no error, got: %s", mkdirError)
	}

	info, err := os.Stat(newDir)
	if err != nil {
		t.Fatalf("Directory was not created: %v", err)
	}

	if !info.IsDir() {
		t.Fatal("Created path is not a directory")
	}
}
