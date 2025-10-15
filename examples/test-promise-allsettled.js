console.log("=== Promise.allSettled() Test ===\n");

// Test 1: Mixed success and failure
console.log("1. Mixed fulfilled and rejected promises:");
Promise.allSettled([
  Promise.resolve(42),
  Promise.reject("Something went wrong"),
  Promise.resolve("success"),
  Promise.reject("Another error")
]).then(results => {
  console.log("  Results:", results);
  results.forEach((result, index) => {
    if (result.status === "fulfilled") {
      console.log(`  [${index}] Fulfilled:`, result.value);
    } else {
      console.log(`  [${index}] Rejected:`, result.reason);
    }
  });
});

// Test 2: All fulfilled
setTimeout(() => {
  console.log("\n2. All promises fulfilled:");
  Promise.allSettled([
    Promise.resolve(1),
    Promise.resolve(2),
    Promise.resolve(3)
  ]).then(results => {
    console.log("  All settled! Count:", results.length);
    results.forEach((r, i) => {
      console.log(`  [${i}]`, r.status, "->", r.value);
    });
  });
}, 500);

// Test 3: All rejected
setTimeout(() => {
  console.log("\n3. All promises rejected:");
  Promise.allSettled([
    Promise.reject("error 1"),
    Promise.reject("error 2"),
    Promise.reject("error 3")
  ]).then(results => {
    console.log("  Still resolves! (never rejects)");
    results.forEach((r, i) => {
      console.log(`  [${i}]`, r.status, "->", r.reason);
    });
  });
}, 1000);

// Test 4: Mixed timing with setTimeout
setTimeout(() => {
  console.log("\n4. Mixed timing:");
  Promise.allSettled([
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise 1 resolving...");
        resolve("fast");
      }, 100);
    }),
    new Promise((resolve, reject) => {
      setTimeout(() => {
        console.log("  Promise 2 rejecting...");
        reject("medium error");
      }, 200);
    }),
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise 3 resolving...");
        resolve("slow");
      }, 300);
    })
  ]).then(results => {
    console.log("\n  All settled after 300ms!");
    results.forEach((r, i) => {
      const statusEmoji = r.status === "fulfilled" ? "✓" : "✗";
      const resultValue = r.status === "fulfilled" ? r.value : r.reason;
      console.log(`  ${statusEmoji} [${i}]`, r.status, "->", resultValue);
    });
  });
}, 1500);

// Test 5: Empty array
setTimeout(() => {
  console.log("\n5. Empty array:");
  Promise.allSettled([]).then(results => {
    console.log("  Results:", results);
    console.log("  Length:", results.length);
  });
}, 2200);

// Test 6: Non-promise values
setTimeout(() => {
  console.log("\n6. Mix of promises and non-promise values:");
  Promise.allSettled([
    42,
    "hello",
    Promise.resolve("promise value"),
    true,
    Promise.reject("promise error")
  ]).then(results => {
    console.log("  All values wrapped as result objects:");
    results.forEach((r, i) => {
      if (r.status === "fulfilled") {
        console.log(`  [${i}] fulfilled:`, r.value);
      } else {
        console.log(`  [${i}] rejected:`, r.reason);
      }
    });
  });
}, 2500);

// Test 7: Comparison with Promise.all()
setTimeout(() => {
  console.log("\n7. Comparison: Promise.all() vs Promise.allSettled()");
  
  console.log("\n  Using Promise.all():");
  Promise.all([
    Promise.resolve("success 1"),
    Promise.reject("FAILURE"),
    Promise.resolve("success 2")
  ])
    .then(results => console.log("  Results:", results))
    .catch(error => console.log("  Caught error:", error));
  
  setTimeout(() => {
    console.log("\n  Using Promise.allSettled():");
    Promise.allSettled([
      Promise.resolve("success 1"),
      Promise.reject("FAILURE"),
      Promise.resolve("success 2")
    ]).then(results => {
      console.log("  Results:");
      results.forEach((r, i) => {
        console.log(`    [${i}] ${r.status}:`, r.status === "fulfilled" ? r.value : r.reason);
      });
    });
  }, 100);
}, 3000);

console.log("\nStarting tests...\n");
