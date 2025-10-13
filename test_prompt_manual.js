console.log("=== Manual Permission Prompt Test ===\n");

console.log("This will prompt for permission to read a file.");
console.log("Try entering 'y', 'n', or 'always'\n");

file.read("/tmp/test.txt", function(err, data) {
  if (err) {
    console.error("Read failed:", err);
  } else {
    console.log("Read succeeded:", data);
  }
  
  console.log("\nTest complete!");
});
