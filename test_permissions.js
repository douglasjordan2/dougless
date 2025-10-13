// Test script for permission system
console.log("Testing permission system...\n");

// Test 1: File read without permission (should fail)
console.log("Test 1: Attempting to read file without permission:");
file.read("/etc/passwd", function(err, data) {
  if (err) {
    console.log("✓ Read correctly denied:", err);
  } else {
    console.log("✗ Read should have been denied!");
  }
});

// Test 2: File write without permission (should fail)
console.log("\nTest 2: Attempting to write file without permission:");
file.write("/tmp/test.txt", "test data", function(err) {
  if (err) {
    console.log("✓ Write correctly denied:", err);
  } else {
    console.log("✗ Write should have been denied!");
  }
});

// Test 3: HTTP request without permission (should fail)
console.log("\nTest 3: Attempting HTTP request without permission:");
http.get("https://api.github.com", function(err, response) {
  if (err) {
    console.log("✓ HTTP correctly denied:", err);
  } else {
    console.log("✗ HTTP should have been denied!");
  }
});
