console.log("=== Promise.any() Test ===\n");

// Test 1: First fulfillment wins
console.log("1. First promise to fulfill wins:");
Promise.any([
  Promise.reject("error 1"),
  Promise.resolve("WINNER!"),
  Promise.resolve("too late")
]).then(result => {
  console.log("  Result:", result);
});

// Test 2: Mixed timing - fastest success wins
setTimeout(() => {
  console.log("\n2. Fastest success wins (with timing):");
  Promise.any([
    new Promise((resolve, reject) => {
      setTimeout(() => reject("slow error"), 200);
    }),
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Fast promise resolving...");
        resolve("fast success");
      }, 100);
    }),
    new Promise((resolve) => {
      setTimeout(() => resolve("slower success"), 300);
    })
  ]).then(result => {
    console.log("  Winner:", result);
  });
}, 500);

// Test 3: All rejected - AggregateError
setTimeout(() => {
  console.log("\n3. All promises rejected (AggregateError):");
  Promise.any([
    Promise.reject("error 1"),
    Promise.reject("error 2"),
    Promise.reject("error 3")
  ]).catch(err => {
    console.log("  Caught AggregateError!");
    console.log("  Name:", err.name);
    console.log("  Message:", err.message);
    console.log("  Errors:", err.errors);
  });
}, 1200);

// Test 4: Empty array
setTimeout(() => {
  console.log("\n4. Empty array:");
  Promise.any([]).catch(err => {
    console.log("  Caught AggregateError!");
    console.log("  Name:", err.name);
    console.log("  Message:", err.message);
  });
}, 1500);

// Test 5: Non-promise values
setTimeout(() => {
  console.log("\n5. Non-promise values (instant success):");
  Promise.any([
    Promise.reject("error"),
    42,  // Non-promise value = instant fulfillment
    Promise.resolve("too late")
  ]).then(result => {
    console.log("  Result:", result);
  });
}, 1800);

// Test 6: One success among many failures
setTimeout(() => {
  console.log("\n6. One success among many failures:");
  Promise.any([
    Promise.reject("fail 1"),
    Promise.reject("fail 2"),
    Promise.resolve("THE ONE SUCCESS"),
    Promise.reject("fail 3"),
    Promise.reject("fail 4")
  ]).then(result => {
    console.log("  Found the success:", result);
  });
}, 2100);

// Test 7: Comparison with Promise.race()
setTimeout(() => {
  console.log("\n7. Comparison: Promise.any() vs Promise.race()");
  
  console.log("\n  Using Promise.race():");
  Promise.race([
    Promise.reject("FAST ERROR"),
    new Promise(resolve => setTimeout(() => resolve("slow success"), 100))
  ])
    .then(result => console.log("  Result:", result))
    .catch(error => console.log("  Caught error:", error));
  
  setTimeout(() => {
    console.log("\n  Using Promise.any():");
    Promise.any([
      Promise.reject("FAST ERROR"),
      new Promise(resolve => setTimeout(() => resolve("slow success"), 100))
    ]).then(result => {
      console.log("  Result:", result);
      console.log("  (Promise.any ignores rejections and waits for success!)");
    });
  }, 50);
}, 2500);

console.log("\nStarting tests...\n");
