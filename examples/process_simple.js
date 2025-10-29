// Simple process module example
console.log('Process Information:');
console.log('- PID:', process.pid);
console.log('- Platform:', process.platform);
console.log('- Arch:', process.arch);
console.log('- Version:', process.version);
console.log('- Current directory:', process.cwd());
console.log('- Arguments:', process.argv);
console.log('- Environment variable USER:', process.env.USER || process.env.USERNAME);

// Exit cleanly
console.log('\nExiting with code 0');
process.exit(0);
