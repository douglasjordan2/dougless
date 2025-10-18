// Crypto Module Demo - HMAC and Timing-Safe Comparison

console.log('=== Crypto Module Demo ===\n');

// 1. HMAC - Hash-based Message Authentication Code
console.log('1. HMAC Examples:');
const secret = 'my-secret-key';
const message = 'Hello, World!';

// SHA256 HMAC (most common)
const hmac256 = crypto.createHmac('sha256', secret);
hmac256.update(message);
console.log('HMAC-SHA256:', hmac256.digest('hex'));

// SHA512 HMAC (more secure)
const hmac512 = crypto.createHmac('sha512', secret);
hmac512.update(message);
console.log('HMAC-SHA512:', hmac512.digest('hex'));

// Base64 encoding (common for APIs)
const hmacBase64 = crypto.createHmac('sha256', secret);
hmacBase64.update(message);
console.log('HMAC (base64):', hmacBase64.digest('base64'));

console.log('\n2. Timing-Safe Comparison:');
// Prevents timing attacks when comparing secrets

const hash1 = crypto.createHash('sha256').update('password123').digest('hex');
const hash2 = crypto.createHash('sha256').update('password123').digest('hex');
const hash3 = crypto.createHash('sha256').update('different').digest('hex');

console.log('Comparing identical hashes:', crypto.timingSafeEqual(hash1, hash2)); // true
console.log('Comparing different hashes:', crypto.timingSafeEqual(hash1, hash3)); // false

// Real-world use case: API signature verification
console.log('\n3. API Signature Verification Example:');
const apiKey = 'sk_test_123456789';
const requestBody = JSON.stringify({ user: 'john', action: 'create' });

// Create signature (server side)
const signature = crypto.createHmac('sha256', apiKey)
  .update(requestBody)
  .digest('hex');

console.log('Request signature:', signature);

// Verify signature (client sends this)
const receivedSignature = crypto.createHmac('sha256', apiKey)
  .update(requestBody)
  .digest('hex');

// Use timingSafeEqual to prevent timing attacks
const isValid = crypto.timingSafeEqual(signature, receivedSignature);
console.log('Signature valid:', isValid);

console.log('\n4. UUID Generation:');
console.log('UUID:', crypto.uuid());
console.log('UUID:', crypto.uuid());

console.log('\n5. Random Bytes:');
console.log('Random hex (16 bytes):', crypto.randomBytes(16, 'hex'));
console.log('Random base64 (32 bytes):', crypto.randomBytes(32, 'base64'));
