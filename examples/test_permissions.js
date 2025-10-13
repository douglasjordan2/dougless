console.log("=== Testing Dougless Permission System ===\n");

// Sequential test - each operation waits for the previous to complete
console.log("1. Attempting to read /tmp/test.txt...");
file.read("/tmp/test.txt", function(err, data) {
  if (err) {
    console.error("   ❌ Read failed:", err);
  } else {
    console.log("   ✓ Read succeeded:", data);
  }
  
  // Only start the next operation after the first completes
  console.log("\n2. Attempting to write to /tmp/dougless_test.txt...");
  file.write("/tmp/dougless_test.txt", "Hello from Dougless!", function(err) {
    if (err) {
      console.error("   ❌ Write failed:", err);
    } else {
      console.log("   ✓ Write succeeded!");
    }
    
    // Only start the third operation after the second completes
    console.log("\n3. Attempting to read directory /tmp...");
    file.readdir("/tmp", function(err, files) {
      if (err) {
        console.error("   ❌ Readdir failed:", err);
      } else {
        console.log("   ✓ Readdir succeeded! Found", files.length, "files");
      }
      
      console.log("\n=== Test Complete ===");
    });
  });
});
