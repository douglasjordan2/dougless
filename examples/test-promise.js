console.log("=== Promise Test ===\n");

// Test 1: Basic promise with setTimeout
console.log("1. Basic Promise with async resolve:");
const p1 = new Promise((resolve, reject) => {
  console.log("  Inside executor");
  setTimeout(() => {
    console.log("  Resolving after 100ms...");
    resolve("Success!");
  }, 100);
});

p1.then(result => {
  console.log("  Resolved:", result);
});

// Test 2: Immediate resolve with Promise.resolve()
setTimeout(() => {
  console.log("\n2. Promise.resolve():");
  Promise.resolve(42).then(value => {
    console.log("  Got value:", value);
  });
}, 200);

// Test 3: Promise chaining
setTimeout(() => {
  console.log("\n3. Promise Chaining:");
  Promise.resolve(5)
    .then(x => {
      console.log("  Step 1:", x);
      return x * 2;
    })
    .then(x => {
      console.log("  Step 2:", x);
      return x + 3;
    })
    .then(x => {
      console.log("  Final:", x);
    });
}, 400);

// Test 4: Error handling with Promise.reject()
setTimeout(() => {
  console.log("\n4. Error Handling with Promise.reject():");
  Promise.reject("Something went wrong")
    .catch(err => {
      console.log("  Caught error:", err);
    });
}, 600);

// Test 5: Error handling with throw in executor
setTimeout(() => {
  console.log("\n5. Promise with rejected executor:");
  new Promise((resolve, reject) => {
    console.log("  About to reject...");
    reject("Manual rejection");
  }).catch(err => {
    console.log("  Caught:", err);
  });
}, 800);

console.log("\nWaiting for promises...");
