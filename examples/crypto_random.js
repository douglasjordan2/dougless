console.log('=== Random Generation ===\n');

// Default (hex)
const hex16 = crypto.random(16);
console.log('16 bytes (hex):', hex16);
console.log('Length:', hex16.length, 'chars'); // Should be 32 (16 bytes × 2 hex chars)

// Base64
const base64_32 = crypto.random(32, 'base64');
console.log('\n32 bytes (base64):', base64_32);

// Raw array
const raw = crypto.random(10, 'raw');
console.log('\n10 bytes (raw):', raw);

// Multiple UUIDs to show uniqueness
console.log('\nUUIDs:');
console.log('1:', crypto.uuid());
console.log('2:', crypto.uuid());
console.log('3:', crypto.uuid());
