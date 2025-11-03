// Test process module
console.log('Testing process module...');

// Test process.env
console.log('Environment variables:');
console.log('  HOME:', process.env.HOME);
console.log('  USER:', process.env.USER);
console.log('  PATH:', process.env.PATH ? 'exists' : 'not set');

// Test process.argv
console.log('\nCommand line arguments:');
console.log('  process.argv:', process.argv);

// Test process.cwd
console.log('\nCurrent working directory:');
console.log('  process.cwd():', process.cwd());

// Test process.platform and process.arch
console.log('\nSystem information:');
console.log('  Platform:', process.platform);
console.log('  Architecture:', process.arch);
console.log('  PID:', process.pid);
console.log('  Version:', process.version);

// Test process.on('exit')
process.on('exit', (code) => {
  console.log(`\nExit handler called with code: ${code}`);
});

// Test process.on('SIGINT') - won't be triggered in this test
process.on('SIGINT', () => {
  console.log('\nReceived SIGINT (Ctrl+C)');
  process.exit(0);
});

console.log('\nProcess module tests complete!');
