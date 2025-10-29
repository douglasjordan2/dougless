// ================================================
// Process Module Demo
// ================================================
// Demonstrates the global process object functionality

console.log('=== Process Module Demo ===\n');

// 1. Command-line arguments
console.log('1. Command-line Arguments:');
console.log('   process.argv:', process.argv);
console.log('   Script name:', process.argv[1] || '(none)');
console.log('   Arguments:', process.argv.slice(2));
console.log();

// 2. Environment variables
console.log('2. Environment Variables:');
console.log('   USER:', process.env.USER || process.env.USERNAME || '(unknown)');
console.log('   HOME:', process.env.HOME || process.env.USERPROFILE || '(unknown)');
console.log('   PATH length:', process.env.PATH ? process.env.PATH.length : 0, 'chars');
console.log('   Total env vars:', Object.keys(process.env).length);
console.log();

// 3. Process information
console.log('3. Process Information:');
console.log('   PID:', process.pid);
console.log('   Platform:', process.platform);
console.log('   Architecture:', process.arch);
console.log('   Dougless version:', process.version);
console.log();

// 4. Working directory
console.log('4. Current Working Directory:');
console.log('   cwd():', process.cwd());
console.log();

// 5. Signal handling
console.log('5. Signal Handling:');
console.log('   Installing SIGINT handler (Ctrl+C)...');

process.on('SIGINT', function(signal) {
  console.log('\n   Caught signal:', signal);
  console.log('   Cleaning up before exit...');
  console.log('   Goodbye!');
  process.exit(0);
});

console.log('   Handler installed. Press Ctrl+C to test.');
console.log();

// 6. Exit handler
console.log('6. Exit Handler:');
process.on('exit', function(code) {
  console.log('   Exit handler called with code:', code);
});

// 7. Timers to keep process alive
console.log('7. Keeping process alive for 10 seconds...');
console.log('   (Press Ctrl+C to exit early)\n');

let countdown = 10;
const timer = setInterval(function() {
  countdown--;
  console.log('   Time remaining:', countdown, 'seconds');
  
  if (countdown <= 0) {
    clearInterval(timer);
    console.log('\n=== Demo Complete ===');
    console.log('Exiting normally with code 0...');
    process.exit(0);
  }
}, 1000);
