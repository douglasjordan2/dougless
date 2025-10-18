# Crypto API Guide

## Overview

Dougless provides a unique global `crypto` API for cryptographic operations. Unlike Node.js which requires `require('crypto')`, the `crypto` object is always available globally.

## Why Global?

**Dougless Philosophy**: Cryptographic operations are fundamental to secure web applications and should be as accessible as `console`. This makes Dougless code cleaner and enables security-first development patterns.

**Comparison:**

```javascript
// Node.js
const crypto = require('crypto');
crypto.createHash('sha256').update('data').digest('hex');

// Dougless
crypto.createHash('sha256').update('data').digest('hex');  // No require!
```

---

## Hash Functions

### `crypto.createHash(algorithm)`

Create a hash object for generating cryptographic hashes.

**Parameters:**
- `algorithm` (string) - Hash algorithm: `'md5'`, `'sha1'`, `'sha256'`, `'sha512'`

**Returns:**
- Hash object with `update()` and `digest()` methods

**Example:**
```javascript
const hash = crypto.createHash('sha256')
    .update('Hello Dougless!')
    .digest('hex');

console.log('SHA-256:', hash);
// SHA-256: a8b2c3d4e5f6...
```

---

### `hash.update(data)`

Add data to be hashed.

**Parameters:**
- `data` (string) - Data to hash

**Returns:**
- Hash object (for chaining)

**Note:** Calling `update()` multiple times will overwrite previous data (not append).

---

### `hash.digest([encoding])`

Generate the hash digest.

**Parameters:**
- `encoding` (string, optional) - Output encoding: `'hex'` (default) or `'base64'`

**Returns:**
- Hash digest as a string

**Examples:**
```javascript
// Hex encoding (default)
const hexHash = crypto.createHash('sha256')
    .update('secure data')
    .digest('hex');

// Base64 encoding
const b64Hash = crypto.createHash('sha256')
    .update('secure data')
    .digest('base64');

console.log('Hex:', hexHash);
console.log('Base64:', b64Hash);
```

---

## HMAC (Message Authentication)

### `crypto.createHmac(algorithm, key)`

Create an HMAC object for message authentication.

**Parameters:**
- `algorithm` (string) - Hash algorithm: `'md5'`, `'sha1'`, `'sha256'`, `'sha512'`
- `key` (string) - Secret key for HMAC

**Returns:**
- HMAC object with `update()` and `digest()` methods

**Example:**
```javascript
const signature = crypto.createHmac('sha256', 'secret-key')
    .update('important message')
    .digest('hex');

console.log('Signature:', signature);
```

---

### HMAC Methods

The HMAC object has the same `update()` and `digest()` methods as hash objects.

**Example - Webhook Signature Verification:**
```javascript
// Server creates signature
const payload = JSON.stringify({ event: 'payment.success', amount: 100 });
const secret = 'webhook-secret-key';

const signature = crypto.createHmac('sha256', secret)
    .update(payload)
    .digest('hex');

// Send signature in header
console.log('X-Signature:', signature);

// Client verifies signature
const computed = crypto.createHmac('sha256', secret)
    .update(payload)
    .digest('hex');

if (crypto.timingSafeEqual(signature, computed)) {
    console.log('✓ Signature valid');
} else {
    console.log('✗ Signature invalid');
}
```

---

## Timing-Safe Comparison

### `crypto.timingSafeEqual(a, b)`

Compare two strings in constant time to prevent timing attacks.

**Parameters:**
- `a` (string) - First string to compare
- `b` (string) - Second string to compare

**Returns:**
- `true` if strings are equal, `false` otherwise

**Why Use This:**
Regular string comparison (`===`) can leak timing information that attackers might exploit. Always use `timingSafeEqual()` when comparing secrets, tokens, or signatures.

**Example:**
```javascript
const storedPassword = 'user-hashed-password';
const providedPassword = 'user-hashed-password';

// ✗ BAD - Timing attack vulnerable
if (storedPassword === providedPassword) {
    console.log('Login success');
}

// ✓ GOOD - Timing attack safe
if (crypto.timingSafeEqual(storedPassword, providedPassword)) {
    console.log('Login success');
}
```

---

## Random Data Generation

### `crypto.random(size, [encoding])`
### `crypto.randomBytes(size, [encoding])` (alias)

Generate cryptographically secure random data.

**Parameters:**
- `size` (number) - Number of bytes to generate (max 65536, i.e., 64KB)
- `encoding` (string, optional) - Output format: `'hex'` (default), `'base64'`, or `'raw'`

**Returns:**
- Random data as a string (hex/base64) or array (raw)

**Examples:**
```javascript
// Hex encoding (default) - 16 bytes = 32 hex characters
const token = crypto.random(16);
console.log('Token:', token);
// Token: a1b2c3d4e5f6g7h8...

// Base64 encoding
const apiKey = crypto.random(32, 'base64');
console.log('API Key:', apiKey);

// Raw bytes as array
const bytes = crypto.random(10, 'raw');
console.log('Bytes:', bytes);
// Bytes: [234, 12, 89, 156, ...]
```

**Use Cases:**
- Session tokens
- API keys
- Cryptographic salts
- Initialization vectors (IVs)
- One-time passwords

---

## UUID Generation

### `crypto.uuid()`

Generate a random UUID (version 4).

**Returns:**
- UUID string in standard format: `xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx`

**Example:**
```javascript
const id = crypto.uuid();
console.log('User ID:', id);
// User ID: 550e8400-e29b-41d4-a716-446655440000

// Generate multiple unique IDs
const ids = [];
for (let i = 0; i < 5; i++) {
    ids.push(crypto.uuid());
}
console.log('IDs:', ids);
```

**Use Cases:**
- Database primary keys
- Session identifiers
- Request tracking IDs
- File upload names

---

## Complete Examples

### Example 1: Password Hashing with Salt

```javascript
// Generate a random salt
const salt = crypto.random(16, 'hex');
console.log('Salt:', salt);

// Hash password with salt
const password = 'user-password-123';
const hash = crypto.createHash('sha256')
    .update(password + salt)
    .digest('hex');

console.log('Hash:', hash);

// Store both salt and hash
const stored = {
    salt: salt,
    hash: hash
};

// Later: verify password
function verifyPassword(inputPassword, stored) {
    const computed = crypto.createHash('sha256')
        .update(inputPassword + stored.salt)
        .digest('hex');
    
    return crypto.timingSafeEqual(computed, stored.hash);
}

console.log('Valid:', verifyPassword('user-password-123', stored)); // true
console.log('Valid:', verifyPassword('wrong-password', stored));     // false
```

---

### Example 2: API Request Signing

```javascript
// API credentials
const apiKey = 'my-api-key';
const apiSecret = 'my-api-secret';

// Create request
const timestamp = Date.now().toString();
const method = 'POST';
const path = '/api/v1/orders';
const body = JSON.stringify({ symbol: 'BTCUSD', qty: 1 });

// Sign request
const message = timestamp + method + path + body;
const signature = crypto.createHmac('sha256', apiSecret)
    .update(message)
    .digest('hex');

console.log('Signature:', signature);

// Make request with signature
http.post('https://api.example.com' + path, 
    JSON.parse(body),
    function(err, response) {
        if (err) {
            console.error('Request failed:', err);
            return;
        }
        console.log('Response:', response.body);
    }
);

// Headers to send:
// X-API-Key: my-api-key
// X-Timestamp: <timestamp>
// X-Signature: <signature>
```

---

### Example 3: Webhook Signature Verification

```javascript
// Server sends webhook with signature
function sendWebhook(url, payload, secret) {
    const body = JSON.stringify(payload);
    const signature = crypto.createHmac('sha256', secret)
        .update(body)
        .digest('hex');
    
    // In real code, send with HTTP headers:
    // X-Webhook-Signature: signature
    
    return { body: body, signature: signature };
}

// Client verifies webhook
function verifyWebhook(body, receivedSignature, secret) {
    const computed = crypto.createHmac('sha256', secret)
        .update(body)
        .digest('hex');
    
    return crypto.timingSafeEqual(receivedSignature, computed);
}

// Usage
const webhookSecret = 'shared-webhook-secret';
const event = { type: 'payment.success', amount: 100 };

const sent = sendWebhook('https://client.com/webhook', event, webhookSecret);

console.log('Signature valid:', 
    verifyWebhook(sent.body, sent.signature, webhookSecret)
); // true

console.log('Wrong signature:', 
    verifyWebhook(sent.body, 'fake-signature', webhookSecret)
); // false
```

---

### Example 4: Secure Token Generation

```javascript
// Generate session token
function createSessionToken() {
    return {
        id: crypto.uuid(),
        token: crypto.random(32, 'base64'),
        createdAt: Date.now()
    };
}

// Generate API key
function createApiKey(userId) {
    const prefix = 'sk_live_';
    const randomPart = crypto.random(24, 'base64')
        .replace(/\+/g, '-')
        .replace(/\//g, '_')
        .replace(/=/g, '');
    
    return prefix + randomPart;
}

const session = createSessionToken();
console.log('Session:', session);

const apiKey = createApiKey('user_123');
console.log('API Key:', apiKey);
```

---

## Supported Algorithms

### Hash Algorithms
- **MD5** - Fast but weak, avoid for security (collision attacks exist)
- **SHA-1** - Legacy, avoid for security (collision attacks exist)
- **SHA-256** - Recommended for most use cases
- **SHA-512** - Stronger, larger output (512 bits vs 256)

### Encoding Options
- **hex** - Hexadecimal (0-9, a-f), 2 characters per byte
- **base64** - Base64 encoding, more compact than hex

---

## Security Best Practices

### ✅ DO:
- Use **SHA-256 or SHA-512** for new applications
- Use **HMAC** for message authentication
- Use **timingSafeEqual()** for comparing secrets
- Use **crypto.random()** for security-sensitive random data
- Store salts alongside hashed passwords
- Generate unique salts for each password

### ❌ DON'T:
- Use MD5 or SHA-1 for security (use for checksums only)
- Compare secrets with `===` (use `timingSafeEqual()`)
- Use `Math.random()` for security (use `crypto.random()`)
- Hash passwords without salt
- Reuse salts across different passwords
- Store passwords in plain text

---

## API Reference Summary

| Method | Purpose | Example |
|--------|---------|---------|
| `createHash(algorithm)` | Create hash object | `crypto.createHash('sha256')` |
| `createHmac(algorithm, key)` | Create HMAC object | `crypto.createHmac('sha256', 'key')` |
| `timingSafeEqual(a, b)` | Safe comparison | `crypto.timingSafeEqual(s1, s2)` |
| `random(size, [encoding])` | Random bytes | `crypto.random(16, 'hex')` |
| `randomBytes(size, [encoding])` | Random bytes (alias) | `crypto.randomBytes(32)` |
| `uuid()` | Generate UUID v4 | `crypto.uuid()` |

---

## Node.js Compatibility

Dougless crypto API is **inspired by** Node.js crypto but simplified:

**Compatible Methods:**
- `createHash()` - Similar API
- `createHmac()` - Similar API
- `randomBytes()` - Similar API
- `timingSafeEqual()` - Same behavior

**Differences:**
- Global access (no `require()`)
- Simplified API (fewer options)
- No streaming support (yet)
- Limited algorithm support (common ones only)

---

*For more examples, see [`examples/crypto_demo.js`](../examples/crypto_demo.js)*
