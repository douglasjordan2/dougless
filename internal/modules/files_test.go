package modules

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/dop251/goja"

	"github.com/douglasjordan2/dougless/internal/event"
	"github.com/douglasjordan2/dougless/internal/permissions"
)

// TestFilesRead_File_Success tests reading an existing file
func TestFilesRead_File_Success(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	// Create temp file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello Dougless!"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Grant read permission
	permissions.GetManager().GrantRead([]string{tmpDir})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	readFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("read"))
	if !ok {
		t.Fatalf("files.read is not a function")
	}

	var (
		cbWg    sync.WaitGroup
		gotErr  goja.Value
		gotData goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		if len(call.Arguments) > 1 {
			gotData = call.Arguments[1]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := readFn(goja.Undefined(), vm.ToValue(testFile), cb)
	if err != nil {
		t.Fatalf("calling files.read failed: %v", err)
	}

	cbWg.Wait()

	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error, got: %v", gotErr)
	}
	if gotData.String() != content {
		t.Fatalf("expected content %q, got %q", content, gotData.String())
	}
}

// TestFilesRead_File_NotExists tests reading a non-existent file (returns null)
func TestFilesRead_File_NotExists(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	tmpDir := t.TempDir()
	nonExistentFile := filepath.Join(tmpDir, "does-not-exist.txt")

	// Grant read permission
	permissions.GetManager().GrantRead([]string{tmpDir})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	readFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("read"))
	if !ok {
		t.Fatalf("files.read is not a function")
	}

	var (
		cbWg    sync.WaitGroup
		gotErr  goja.Value
		gotData goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		if len(call.Arguments) > 1 {
			gotData = call.Arguments[1]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := readFn(goja.Undefined(), vm.ToValue(nonExistentFile), cb)
	if err != nil {
		t.Fatalf("calling files.read failed: %v", err)
	}

	cbWg.Wait()

	// Should return null (not error) when file doesn't exist
	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error, got: %v", gotErr)
	}
	if !goja.IsNull(gotData) {
		t.Fatalf("expected null data for non-existent file, got: %v", gotData)
	}
}

// TestFilesRead_Directory_Success tests listing directory contents
func TestFilesRead_Directory_Success(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	// Create temp directory with files
	tmpDir := t.TempDir()
	files := []string{"a.txt", "b.txt", "c.txt"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Grant read permission
	permissions.GetManager().GrantRead([]string{tmpDir})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	readFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("read"))
	if !ok {
		t.Fatalf("files.read is not a function")
	}

	var (
		cbWg    sync.WaitGroup
		gotErr  goja.Value
		gotData goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		if len(call.Arguments) > 1 {
			gotData = call.Arguments[1]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	// Note the trailing slash to indicate directory
	_, err := readFn(goja.Undefined(), vm.ToValue(tmpDir+"/"), cb)
	if err != nil {
		t.Fatalf("calling files.read failed: %v", err)
	}

	cbWg.Wait()

	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error, got: %v", gotErr)
	}

	// Should return array of filenames
	arr := gotData.Export().([]string)
	if len(arr) != 3 {
		t.Fatalf("expected 3 files, got %d", len(arr))
	}

	// Check that all expected files are present
	fileMap := make(map[string]bool)
	for _, f := range arr {
		fileMap[f] = true
	}
	for _, expected := range files {
		if !fileMap[expected] {
			t.Fatalf("expected file %q not found in results", expected)
		}
	}
}

// TestFilesWrite_File_Success tests writing a file
func TestFilesWrite_File_Success(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "output.txt")
	content := "Hello World!"

	// Grant write permission
	permissions.GetManager().GrantWrite([]string{tmpDir})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	writeFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("write"))
	if !ok {
		t.Fatalf("files.write is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := writeFn(goja.Undefined(), vm.ToValue(testFile), vm.ToValue(content), cb)
	if err != nil {
		t.Fatalf("calling files.write failed: %v", err)
	}

	cbWg.Wait()

	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error, got: %v", gotErr)
	}

	// Verify file was created with correct content
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read created file: %v", err)
	}
	if string(data) != content {
		t.Fatalf("expected content %q, got %q", content, string(data))
	}
}

// TestFilesWrite_File_AutoCreateParentDirs tests auto-creation of parent directories
func TestFilesWrite_File_AutoCreateParentDirs(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	tmpDir := t.TempDir()
	nestedFile := filepath.Join(tmpDir, "a", "b", "c", "test.txt")
	content := "Nested!"

	// Grant write permission
	permissions.GetManager().GrantWrite([]string{tmpDir})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	writeFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("write"))
	if !ok {
		t.Fatalf("files.write is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := writeFn(goja.Undefined(), vm.ToValue(nestedFile), vm.ToValue(content), cb)
	if err != nil {
		t.Fatalf("calling files.write failed: %v", err)
	}

	cbWg.Wait()

	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error, got: %v", gotErr)
	}

	// Verify file was created and parent directories exist
	data, err := os.ReadFile(nestedFile)
	if err != nil {
		t.Fatalf("failed to read created file: %v", err)
	}
	if string(data) != content {
		t.Fatalf("expected content %q, got %q", content, string(data))
	}

	// Verify parent directories were created
	parentDir := filepath.Join(tmpDir, "a", "b", "c")
	stat, err := os.Stat(parentDir)
	if err != nil {
		t.Fatalf("parent directory was not created: %v", err)
	}
	if !stat.IsDir() {
		t.Fatalf("expected parent path to be a directory")
	}
}

// TestFilesWrite_Directory_Success tests creating a directory
func TestFilesWrite_Directory_Success(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "newdir") + "/"

	// Grant write permission
	permissions.GetManager().GrantWrite([]string{tmpDir})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	writeFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("write"))
	if !ok {
		t.Fatalf("files.write is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	// Note: only 2 arguments (path with trailing slash, callback)
	_, err := writeFn(goja.Undefined(), vm.ToValue(newDir), cb)
	if err != nil {
		t.Fatalf("calling files.write failed: %v", err)
	}

	cbWg.Wait()

	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error, got: %v", gotErr)
	}

	// Verify directory was created
	stat, err := os.Stat(filepath.Join(tmpDir, "newdir"))
	if err != nil {
		t.Fatalf("directory was not created: %v", err)
	}
	if !stat.IsDir() {
		t.Fatalf("expected path to be a directory")
	}
}

// TestFilesRm_File_Success tests removing a file
func TestFilesRm_File_Success(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	// Create temp file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "remove-me.txt")
	if err := os.WriteFile(testFile, []byte("delete this"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Grant write permission
	permissions.GetManager().GrantWrite([]string{tmpDir})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	rmFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("rm"))
	if !ok {
		t.Fatalf("files.rm is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := rmFn(goja.Undefined(), vm.ToValue(testFile), cb)
	if err != nil {
		t.Fatalf("calling files.rm failed: %v", err)
	}

	cbWg.Wait()

	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error, got: %v", gotErr)
	}

	// Verify file was deleted
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Fatalf("file should have been deleted")
	}
}

// TestFilesRm_Directory_Recursive tests removing a directory with contents
func TestFilesRm_Directory_Recursive(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	// Create directory with nested content
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "deleteme")
	os.Mkdir(testDir, 0755)
	os.Mkdir(filepath.Join(testDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(testDir, "file1.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(testDir, "subdir", "file2.txt"), []byte("content"), 0644)

	// Grant write permission
	permissions.GetManager().GrantWrite([]string{tmpDir})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	rmFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("rm"))
	if !ok {
		t.Fatalf("files.rm is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := rmFn(goja.Undefined(), vm.ToValue(testDir), cb)
	if err != nil {
		t.Fatalf("calling files.rm failed: %v", err)
	}

	cbWg.Wait()

	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error, got: %v", gotErr)
	}

	// Verify directory and all contents were deleted
	if _, err := os.Stat(testDir); !os.IsNotExist(err) {
		t.Fatalf("directory should have been deleted recursively")
	}
}

// TestFilesRm_NonExistent_Idempotent tests that removing non-existent path doesn't error
func TestFilesRm_NonExistent_Idempotent(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	tmpDir := t.TempDir()
	nonExistentPath := filepath.Join(tmpDir, "does-not-exist.txt")

	// Grant write permission
	permissions.GetManager().GrantWrite([]string{tmpDir})

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	rmFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("rm"))
	if !ok {
		t.Fatalf("files.rm is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := rmFn(goja.Undefined(), vm.ToValue(nonExistentPath), cb)
	if err != nil {
		t.Fatalf("calling files.rm failed: %v", err)
	}

	cbWg.Wait()

	// Should succeed even though file doesn't exist (idempotent)
	if !goja.IsNull(gotErr) {
		t.Fatalf("expected null error for non-existent path, got: %v", gotErr)
	}
}

// TestFilesRead_PermissionDenied tests that read fails without permission
func TestFilesRead_PermissionDenied(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	// Create temp file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("secret"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// DO NOT grant permission

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	readFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("read"))
	if !ok {
		t.Fatalf("files.read is not a function")
	}

	var (
		cbWg    sync.WaitGroup
		gotErr  goja.Value
		gotData goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		if len(call.Arguments) > 1 {
			gotData = call.Arguments[1]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := readFn(goja.Undefined(), vm.ToValue(testFile), cb)
	if err != nil {
		t.Fatalf("calling files.read failed: %v", err)
	}

	cbWg.Wait()

	// Should get an error
	if goja.IsNull(gotErr) || goja.IsUndefined(gotErr) {
		t.Fatalf("expected permission error, got null/undefined")
	}

	// Data should be undefined
	if !goja.IsUndefined(gotData) {
		t.Fatalf("expected undefined data on error, got: %v", gotData)
	}
}

// TestFilesWrite_PermissionDenied tests that write fails without permission
func TestFilesWrite_PermissionDenied(t *testing.T) {
	cleanupPerms := withFreshPermissions(t)
	defer cleanupPerms()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "output.txt")

	// DO NOT grant permission

	loop := event.NewLoop()
	go loop.Run()
	defer func() { loop.Stop(); loop.Wait() }()

	vm := goja.New()
	filesMod := NewFiles(loop)
	filesObj := filesMod.Export(vm)

	writeFn, ok := goja.AssertFunction(filesObj.ToObject(vm).Get("write"))
	if !ok {
		t.Fatalf("files.write is not a function")
	}

	var (
		cbWg   sync.WaitGroup
		gotErr goja.Value
	)
	cbWg.Add(1)

	cb := vm.ToValue(func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			gotErr = call.Arguments[0]
		}
		cbWg.Done()
		return goja.Undefined()
	})

	_, err := writeFn(goja.Undefined(), vm.ToValue(testFile), vm.ToValue("content"), cb)
	if err != nil {
		t.Fatalf("calling files.write failed: %v", err)
	}

	cbWg.Wait()

	// Should get an error
	if goja.IsNull(gotErr) || goja.IsUndefined(gotErr) {
		t.Fatalf("expected permission error, got null/undefined")
	}

	// File should not have been created
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Fatalf("file should not have been created without permission")
	}
}
