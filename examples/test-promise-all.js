console.log("=== Promise.all() Test ===\n");

// Test 1: All promises resolve
console.log("1. All promises resolve:");
Promise.all([
  Promise.resolve(1),
  Promise.resolve(2),
  Promise.resolve(3)
]).then(results => {
  console.log("  Results:", results);
  console.log("  First:", results[0]);
  console.log("  Second:", results[1]);
  console.log("  Third:", results[2]);
});

// Test 2: Mixed timing
setTimeout(() => {
  console.log("\n2. Promises with different timing:");
  Promise.all([
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise 1 resolving...");
        resolve("fast");
      }, 100);
    }),
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise 2 resolving...");
        resolve("slow");
      }, 300);
    }),
    Promise.resolve("immediate")
  ]).then(results => {
    console.log("  All done! Results:", results);
  });
}, 500);

// Test 3: One promise rejects
setTimeout(() => {
  console.log("\n3. One promise rejects:");
  Promise.all([
    Promise.resolve("success 1"),
    Promise.reject("FAILURE"),
    Promise.resolve("success 2")
  ]).then(results => {
    console.log("  This shouldn't print!");
  }).catch(error => {
    console.log("  Caught error:", error);
  });
}, 1500);

// Test 4: Empty array
setTimeout(() => {
  console.log("\n4. Empty array:");
  Promise.all([]).then(results => {
    console.log("  Results:", results);
    console.log("  Length:", results.length);
  });
}, 2000);

// Test 5: Non-promise values
setTimeout(() => {
  console.log("\n5. Mix of promises and values:");
  Promise.all([
    42,
    Promise.resolve("hello"),
    true,
    Promise.resolve(100)
  ]).then(results => {
    console.log("  Results:", results);
  });
}, 2500);

console.log("\nStarting tests...\n");
