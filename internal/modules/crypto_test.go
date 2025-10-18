package modules

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"regexp"
	"testing"

	"github.com/dop251/goja"
)

func setupCryptoTest() (*goja.Runtime, *Crypto) {
	vm := goja.New()
	crypto := NewCrypto()
	vm.Set("crypto", crypto.Export(vm))
	return vm, crypto
}

// Hash Function Tests

func TestCreateHashMD5(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHash('md5')
			.update('hello world')
			.digest('hex');
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	hash := md5.Sum([]byte("hello world"))
	expected := hex.EncodeToString(hash[:])

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCreateHashSHA1(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHash('sha1')
			.update('test data')
			.digest('hex');
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	hash := sha1.Sum([]byte("test data"))
	expected := hex.EncodeToString(hash[:])

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCreateHashSHA256(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHash('sha256')
			.update('secure data')
			.digest('hex');
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	hash := sha256.Sum256([]byte("secure data"))
	expected := hex.EncodeToString(hash[:])

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCreateHashSHA512(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHash('sha512')
			.update('very secure')
			.digest('hex');
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	hash := sha512.Sum512([]byte("very secure"))
	expected := hex.EncodeToString(hash[:])

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCreateHashBase64(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHash('sha256')
			.update('hello')
			.digest('base64');
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	hash := sha256.Sum256([]byte("hello"))
	expected := base64.StdEncoding.EncodeToString(hash[:])

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCreateHashChaining(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		var h = crypto.createHash('sha256');
		h.update('first');
		h.update('second'); // Should overwrite
		h.digest('hex');
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	// Last update wins in current implementation
	hash := sha256.Sum256([]byte("second"))
	expected := hex.EncodeToString(hash[:])

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCreateHashInvalidAlgorithm(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHash('invalid')
			.update('test')
			.digest('hex');
	`

	_, err := vm.RunString(script)
	if err == nil {
		t.Error("Expected error for invalid algorithm, got none")
	}
}

func TestCreateHashMissingArgument(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.createHash();`

	_, err := vm.RunString(script)
	if err == nil {
		t.Error("Expected error for missing argument, got none")
	}
}

func TestCreateHashInvalidEncoding(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHash('sha256')
			.update('test')
			.digest('invalid');
	`

	_, err := vm.RunString(script)
	if err == nil {
		t.Error("Expected error for invalid encoding, got none")
	}
}

// HMAC Tests

func TestCreateHmacSHA256(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHmac('sha256', 'secret-key')
			.update('message to authenticate')
			.digest('hex');
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	mac := hmac.New(sha256.New, []byte("secret-key"))
	mac.Write([]byte("message to authenticate"))
	expected := hex.EncodeToString(mac.Sum(nil))

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCreateHmacMD5(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHmac('md5', 'key')
			.update('data')
			.digest('hex');
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	mac := hmac.New(md5.New, []byte("key"))
	mac.Write([]byte("data"))
	expected := hex.EncodeToString(mac.Sum(nil))

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCreateHmacSHA512Base64(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHmac('sha512', 'my-secret')
			.update('important data')
			.digest('base64');
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	mac := hmac.New(sha512.New, []byte("my-secret"))
	mac.Write([]byte("important data"))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if result.String() != expected {
		t.Errorf("Expected %s, got %s", expected, result.String())
	}
}

func TestCreateHmacMissingArguments(t *testing.T) {
	vm, _ := setupCryptoTest()

	tests := []string{
		`crypto.createHmac();`,
		`crypto.createHmac('sha256');`,
	}

	for _, script := range tests {
		_, err := vm.RunString(script)
		if err == nil {
			t.Errorf("Expected error for script: %s", script)
		}
	}
}

func TestCreateHmacInvalidAlgorithm(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		crypto.createHmac('invalid', 'key')
			.update('data')
			.digest('hex');
	`

	_, err := vm.RunString(script)
	if err == nil {
		t.Error("Expected error for invalid algorithm, got none")
	}
}

// Timing Safe Equal Tests

func TestTimingSafeEqualTrue(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.timingSafeEqual('password123', 'password123');`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if !result.ToBoolean() {
		t.Error("Expected true for equal strings")
	}
}

func TestTimingSafeEqualFalse(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.timingSafeEqual('password123', 'password124');`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if result.ToBoolean() {
		t.Error("Expected false for different strings")
	}
}

func TestTimingSafeEqualDifferentLengths(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.timingSafeEqual('short', 'longer string');`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if result.ToBoolean() {
		t.Error("Expected false for strings of different lengths")
	}
}

func TestTimingSafeEqualEmpty(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.timingSafeEqual('', '');`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if !result.ToBoolean() {
		t.Error("Expected true for empty strings")
	}
}

func TestTimingSafeEqualMissingArguments(t *testing.T) {
	vm, _ := setupCryptoTest()

	tests := []string{
		`crypto.timingSafeEqual();`,
		`crypto.timingSafeEqual('one');`,
	}

	for _, script := range tests {
		_, err := vm.RunString(script)
		if err == nil {
			t.Errorf("Expected error for script: %s", script)
		}
	}
}

// UUID Tests

func TestUUID(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.uuid();`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	uuid := result.String()

	// Check UUID v4 format: 8-4-4-4-12 hex digits
	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !uuidRegex.MatchString(uuid) {
		t.Errorf("Invalid UUID format: %s", uuid)
	}
}

func TestUUIDUniqueness(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		var uuid1 = crypto.uuid();
		var uuid2 = crypto.uuid();
		uuid1 !== uuid2;
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if !result.ToBoolean() {
		t.Error("Expected UUIDs to be unique")
	}
}

// Random Bytes Tests

func TestRandomBytesHex(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.random(16, 'hex');`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	hexString := result.String()
	if len(hexString) != 32 { // 16 bytes = 32 hex characters
		t.Errorf("Expected 32 hex characters, got %d", len(hexString))
	}

	// Verify it's valid hex
	_, err = hex.DecodeString(hexString)
	if err != nil {
		t.Errorf("Invalid hex string: %v", err)
	}
}

func TestRandomBytesBase64(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.random(16, 'base64');`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	b64String := result.String()

	// Verify it's valid base64
	_, err = base64.StdEncoding.DecodeString(b64String)
	if err != nil {
		t.Errorf("Invalid base64 string: %v", err)
	}
}

func TestRandomBytesRaw(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		var bytes = crypto.random(10, 'raw');
		bytes.length;
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if result.ToInteger() != 10 {
		t.Errorf("Expected array of length 10, got %d", result.ToInteger())
	}
}

func TestRandomBytesDefaultHex(t *testing.T) {
	vm, _ := setupCryptoTest()

	// Default encoding should be hex
	script := `crypto.random(8);`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	hexString := result.String()
	if len(hexString) != 16 { // 8 bytes = 16 hex characters
		t.Errorf("Expected 16 hex characters, got %d", len(hexString))
	}
}

func TestRandomBytesAlias(t *testing.T) {
	vm, _ := setupCryptoTest()

	// randomBytes should be an alias for random
	script := `crypto.randomBytes(16, 'hex');`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	hexString := result.String()
	if len(hexString) != 32 {
		t.Errorf("Expected 32 hex characters, got %d", len(hexString))
	}
}

func TestRandomBytesUniqueness(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		var r1 = crypto.random(16);
		var r2 = crypto.random(16);
		r1 !== r2;
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if !result.ToBoolean() {
		t.Error("Expected random bytes to be unique")
	}
}

func TestRandomBytesNegativeSize(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.random(-1);`

	_, err := vm.RunString(script)
	if err == nil {
		t.Error("Expected error for negative size")
	}
}

func TestRandomBytesTooLarge(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.random(70000);` // Over 64kb limit

	_, err := vm.RunString(script)
	if err == nil {
		t.Error("Expected error for size over 64kb")
	}
}

func TestRandomBytesMissingArgument(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.random();`

	_, err := vm.RunString(script)
	if err == nil {
		t.Error("Expected error for missing size argument")
	}
}

func TestRandomBytesInvalidEncoding(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `crypto.random(16, 'invalid');`

	_, err := vm.RunString(script)
	if err == nil {
		t.Error("Expected error for invalid encoding")
	}
}

// Integration Tests

func TestCryptoAPIExists(t *testing.T) {
	vm, _ := setupCryptoTest()

	methods := []string{
		"createHash",
		"createHmac",
		"timingSafeEqual",
		"random",
		"randomBytes",
		"uuid",
	}

	for _, method := range methods {
		script := `typeof crypto.` + method
		result, err := vm.RunString(script)
		if err != nil {
			t.Fatalf("Failed to check method %s: %v", method, err)
		}

		if result.String() != "function" {
			t.Errorf("Expected crypto.%s to be a function, got %s", method, result.String())
		}
	}
}

func TestRealWorldHmacVerification(t *testing.T) {
	vm, _ := setupCryptoTest()

	// Simulate webhook signature verification
	script := `
		var secret = 'webhook-secret';
		var payload = '{"event":"payment.success"}';
		
		// Server creates signature
		var signature = crypto.createHmac('sha256', secret)
			.update(payload)
			.digest('hex');
		
		// Client verifies signature
		var computed = crypto.createHmac('sha256', secret)
			.update(payload)
			.digest('hex');
		
		crypto.timingSafeEqual(signature, computed);
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if !result.ToBoolean() {
		t.Error("HMAC verification failed for identical signatures")
	}
}

func TestPasswordHashing(t *testing.T) {
	vm, _ := setupCryptoTest()

	script := `
		var password = 'user-password';
		var salt = crypto.random(16, 'hex');
		
		// Hash password with salt
		var hash = crypto.createHash('sha256')
			.update(password + salt)
			.digest('hex');
		
		// Verify by recomputing
		var verify = crypto.createHash('sha256')
			.update(password + salt)
			.digest('hex');
		
		crypto.timingSafeEqual(hash, verify);
	`

	result, err := vm.RunString(script)
	if err != nil {
		t.Fatalf("Script execution failed: %v", err)
	}

	if !result.ToBoolean() {
		t.Error("Password verification failed")
	}
}
