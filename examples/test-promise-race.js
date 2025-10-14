console.log("=== Promise.race() Test ===\n");

// Test 1: First promise wins (immediate)
console.log("1. Immediate winner:");
Promise.race([
  Promise.resolve("fast"),
  Promise.resolve("slow")
]).then(result => {
  console.log("  Winner:", result);
});

// Test 2: Different timing - fastest wins
setTimeout(() => {
  console.log("\n2. Different timing (fastest wins):");
  Promise.race([
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise 1 (100ms) resolving...");
        resolve("100ms winner");
      }, 100);
    }),
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise 2 (300ms) resolving...");
        resolve("300ms - too slow");
      }, 300);
    }),
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise 3 (500ms) resolving...");
        resolve("500ms - way too slow");
      }, 500);
    })
  ]).then(result => {
    console.log("  Race winner:", result);
  });
}, 500);

// Test 3: First to reject wins (rejection wins the race)
setTimeout(() => {
  console.log("\n3. First to reject:");
  Promise.race([
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Slow promise resolving...");
        resolve("I'm slow");
      }, 300);
    }),
    new Promise((resolve, reject) => {
      setTimeout(() => {
        console.log("  Fast promise rejecting...");
        reject("FAST FAILURE");
      }, 100);
    })
  ]).then(result => {
    console.log("  This shouldn't print!");
  }).catch(error => {
    console.log("  Caught rejection:", error);
  });
}, 2000);

// Test 4: Non-promise value (wins immediately)
setTimeout(() => {
  console.log("\n4. Non-promise value (instant winner):");
  Promise.race([
    "instant winner",
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  This promise is too slow");
        resolve("too late");
      }, 100);
    })
  ]).then(result => {
    console.log("  Winner:", result);
  });
}, 3000);

// Test 5: Empty array (forever pending - won't resolve)
setTimeout(() => {
  console.log("\n5. Empty array (forever pending):");
  console.log("  Racing empty array...");
  Promise.race([]).then(result => {
    console.log("  This will never print!");
  });
  setTimeout(() => {
    console.log("  Still waiting... (as expected, never resolves)");
  }, 200);
}, 3500);

// Test 6: Timeout pattern (common use case)
setTimeout(() => {
  console.log("\n6. Timeout pattern (race against time):");
  
  var slowOperation = new Promise((resolve) => {
    setTimeout(() => {
      resolve("Operation completed");
    }, 500);
  });
  
  var timeout = new Promise((resolve, reject) => {
    setTimeout(() => {
      reject("TIMEOUT: Operation took too long");
    }, 200);
  });
  
  Promise.race([slowOperation, timeout])
    .then(result => {
      console.log("  Success:", result);
    })
    .catch(error => {
      console.log("  Failed:", error);
    });
}, 4200);

// Test 7: All promises resolve at similar times
setTimeout(() => {
  console.log("\n7. Close race (all ~100ms, first wins):");
  Promise.race([
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise A finishing");
        resolve("A");
      }, 100);
    }),
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise B finishing");
        resolve("B");
      }, 100);
    }),
    new Promise((resolve) => {
      setTimeout(() => {
        console.log("  Promise C finishing");
        resolve("C");
      }, 100);
    })
  ]).then(result => {
    console.log("  Winner:", result);
  });
}, 5000);

console.log("\nStarting tests...\n");
