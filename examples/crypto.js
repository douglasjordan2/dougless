// Test crypto module
console.log('Testing crypto module...');

const testData = 'Hello, Dougless!';

// Test MD5
const md5Hash = crypto.createHash('md5').update(testData).digest('hex');
console.log('MD5 hash:', md5Hash);

// Test SHA1
const sha1Hash = crypto.createHash('sha1').update(testData).digest('hex');
console.log('SHA1 hash:', sha1Hash);

// Test SHA256
const sha256Hash = crypto.createHash('sha256').update(testData).digest('hex');
console.log('SHA256 hash:', sha256Hash);

// Test SHA512
const sha512Hash = crypto.createHash('sha512').update(testData).digest('hex');
console.log('SHA512 hash:', sha512Hash);

// Test with base64 encoding
const base64Hash = crypto.createHash('sha256').update(testData).digest('base64');
console.log('SHA256 (base64):', base64Hash);

// Test HMAC
const hmacHash = crypto.createHmac('sha256', 'secret-key').update(testData).digest('hex');
console.log('HMAC-SHA256:', hmacHash);

// Test UUID
const id = crypto.uuid();
console.log('UUID:', id);

// Test random bytes
const randomHex = crypto.random(16);
console.log('Random (16 bytes):', randomHex);

// Test with empty string
console.log('\nTesting with empty string:');
const emptyMd5 = crypto.createHash('md5').update('').digest('hex');
console.log('MD5:', emptyMd5);

console.log('\nCrypto module tests complete!');
