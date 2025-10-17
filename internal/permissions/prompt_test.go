package permissions

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockPrompter is a test implementation of the Prompter interface
type MockPrompter struct {
	Response PromptResponse
	Err      error
	Called   int
}

func (m *MockPrompter) Prompt(ctx context.Context, desc PermissionDescriptor) (PromptResponse, error) {
	m.Called++

	// Check if context is cancelled
	select {
	case <-ctx.Done():
		return PromptResponse{Granted: false}, ctx.Err()
	default:
	}

	if m.Err != nil {
		return PromptResponse{Granted: false}, m.Err
	}
	return m.Response, nil
}

func TestMockPrompter_GrantPermission(t *testing.T) {
	mock := &MockPrompter{
		Response: PromptResponse{Granted: true, SaveToConfig: false},
	}

	desc := PermissionDescriptor{
		Name:     PermissionRead,
		Resource: "/tmp/test.txt",
	}

	response, err := mock.Prompt(context.Background(), desc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !response.Granted {
		t.Error("expected permission to be granted")
	}

	if response.SaveToConfig {
		t.Error("expected SaveToConfig to be false")
	}

	if mock.Called != 1 {
		t.Errorf("expected Prompt to be called once, got %d", mock.Called)
	}
}

func TestMockPrompter_DenyPermission(t *testing.T) {
	mock := &MockPrompter{
		Response: PromptResponse{Granted: false, SaveToConfig: false},
	}

	desc := PermissionDescriptor{
		Name:     PermissionNet,
		Resource: "api.example.com",
	}

	response, err := mock.Prompt(context.Background(), desc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Granted {
		t.Error("expected permission to be denied")
	}

	if response.SaveToConfig {
		t.Error("SaveToConfig should not be set when permission denied")
	}
}

func TestMockPrompter_GrantAndSaveToConfig(t *testing.T) {
	mock := &MockPrompter{
		Response: PromptResponse{Granted: true, SaveToConfig: true},
	}

	desc := PermissionDescriptor{
		Name:     PermissionWrite,
		Resource: "/tmp/output",
	}

	response, err := mock.Prompt(context.Background(), desc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !response.Granted {
		t.Error("expected permission to be granted")
	}

	if !response.SaveToConfig {
		t.Error("expected SaveToConfig to be true")
	}
}

func TestMockPrompter_Error(t *testing.T) {
	expectedErr := errors.New("mock error")
	mock := &MockPrompter{
		Err: expectedErr,
	}

	desc := PermissionDescriptor{
		Name:     PermissionRead,
		Resource: "/etc/passwd",
	}

	response, err := mock.Prompt(context.Background(), desc)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if response.Granted {
		t.Error("expected permission to be denied on error")
	}
}

func TestMockPrompter_ContextCancellation(t *testing.T) {
	mock := &MockPrompter{
		Response: PromptResponse{Granted: true, SaveToConfig: false},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	desc := PermissionDescriptor{
		Name:     PermissionNet,
		Resource: "localhost:3000",
	}

	response, err := mock.Prompt(ctx, desc)
	if err == nil {
		t.Fatal("expected context cancellation error")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	if response.Granted {
		t.Error("expected permission to be denied on context cancellation")
	}
}

func TestMockPrompter_ContextTimeout(t *testing.T) {
	mock := &MockPrompter{
		Response: PromptResponse{Granted: true, SaveToConfig: false},
	}

	// Create context that expires immediately
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Give it time to expire
	time.Sleep(10 * time.Millisecond)

	desc := PermissionDescriptor{
		Name:     PermissionRun,
		Resource: "git",
	}

	response, err := mock.Prompt(ctx, desc)
	if err == nil {
		t.Fatal("expected timeout error")
	}

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected context.DeadlineExceeded error, got %v", err)
	}

	if response.Granted {
		t.Error("expected permission to be denied on timeout")
	}
}

func TestMockPrompter_MultiplePrompts(t *testing.T) {
	mock := &MockPrompter{
		Response: PromptResponse{Granted: true, SaveToConfig: false},
	}

	// First prompt
	desc1 := PermissionDescriptor{Name: PermissionRead, Resource: "/tmp/file1.txt"}
	_, err := mock.Prompt(context.Background(), desc1)
	if err != nil {
		t.Fatalf("first prompt failed: %v", err)
	}

	// Second prompt
	desc2 := PermissionDescriptor{Name: PermissionRead, Resource: "/tmp/file2.txt"}
	_, err = mock.Prompt(context.Background(), desc2)
	if err != nil {
		t.Fatalf("second prompt failed: %v", err)
	}

	if mock.Called != 2 {
		t.Errorf("expected Prompt to be called twice, got %d", mock.Called)
	}
}

func TestPermissionDescriptor_String(t *testing.T) {
	tests := []struct {
		name     string
		desc     PermissionDescriptor
		expected string
	}{
		{
			name:     "read with resource",
			desc:     PermissionDescriptor{Name: PermissionRead, Resource: "/tmp/test.txt"},
			expected: "read access to '/tmp/test.txt'",
		},
		{
			name:     "write with resource",
			desc:     PermissionDescriptor{Name: PermissionWrite, Resource: "/var/log/app.log"},
			expected: "write access to '/var/log/app.log'",
		},
		{
			name:     "net with resource",
			desc:     PermissionDescriptor{Name: PermissionNet, Resource: "api.example.com"},
			expected: "net access to 'api.example.com'",
		},
		{
			name:     "env with resource",
			desc:     PermissionDescriptor{Name: PermissionEnv, Resource: "API_KEY"},
			expected: "env access to 'API_KEY'",
		},
		{
			name:     "run with resource",
			desc:     PermissionDescriptor{Name: PermissionRun, Resource: "git"},
			expected: "run access to 'git'",
		},
		{
			name:     "read without resource",
			desc:     PermissionDescriptor{Name: PermissionRead, Resource: ""},
			expected: "read access",
		},
		{
			name:     "net without resource",
			desc:     PermissionDescriptor{Name: PermissionNet, Resource: ""},
			expected: "net access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.desc.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestStdioPrompter_Creation(t *testing.T) {
	prompter := NewStdioPrompter()
	if prompter == nil {
		t.Fatal("NewStdioPrompter returned nil")
	}
}

// Integration test with Manager using mock prompter
func TestManagerWithMockPrompter(t *testing.T) {
	manager := NewManager()
	manager.SetPromptMode(true)

	mock := &MockPrompter{
		Response: PromptResponse{Granted: true, SaveToConfig: false},
	}
	manager.SetPrompter(mock)

	// First check - should prompt
	desc := PermissionDescriptor{Name: PermissionRead, Resource: "/tmp/test.txt"}
	granted := manager.CheckWithPrompt(context.Background(), desc.Name, desc.Resource)

	if !granted {
		t.Error("expected permission to be granted")
	}

	if mock.Called != 1 {
		t.Errorf("expected mock to be called once, got %d", mock.Called)
	}

	// Second check for same resource - should use cache, not prompt again
	granted = manager.CheckWithPrompt(context.Background(), desc.Name, desc.Resource)

	if !granted {
		t.Error("expected cached permission to be granted")
	}

	if mock.Called != 1 {
		t.Errorf("expected mock to still be called once (cached), got %d", mock.Called)
	}
}

func TestManagerWithMockPrompter_Deny(t *testing.T) {
	manager := NewManager()
	manager.SetPromptMode(true)

	mock := &MockPrompter{
		Response: PromptResponse{Granted: false, SaveToConfig: false},
	}
	manager.SetPrompter(mock)

	desc := PermissionDescriptor{Name: PermissionNet, Resource: "evil.com"}
	granted := manager.CheckWithPrompt(context.Background(), desc.Name, desc.Resource)

	if granted {
		t.Error("expected permission to be denied")
	}

	if mock.Called != 1 {
		t.Errorf("expected mock to be called once, got %d", mock.Called)
	}

	// Second check - should use cached denial
	granted = manager.CheckWithPrompt(context.Background(), desc.Name, desc.Resource)

	if granted {
		t.Error("expected cached denial")
	}

	if mock.Called != 1 {
		t.Errorf("expected mock to still be called once (cached), got %d", mock.Called)
	}
}

func TestManagerWithMockPrompter_ClearCache(t *testing.T) {
	manager := NewManager()
	manager.SetPromptMode(true)

	mock := &MockPrompter{
		Response: PromptResponse{Granted: true, SaveToConfig: false},
	}
	manager.SetPrompter(mock)

	desc := PermissionDescriptor{Name: PermissionWrite, Resource: "/tmp/output.txt"}

	// First prompt
	manager.CheckWithPrompt(context.Background(), desc.Name, desc.Resource)
	if mock.Called != 1 {
		t.Errorf("expected 1 call, got %d", mock.Called)
	}

	// Clear cache
	manager.ClearPromptCache()

	// Second prompt - should prompt again after cache clear
	manager.CheckWithPrompt(context.Background(), desc.Name, desc.Resource)
	if mock.Called != 2 {
		t.Errorf("expected 2 calls after cache clear, got %d", mock.Called)
	}
}
