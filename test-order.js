console.log("=== Promise.race ===");
Promise.race([
  Promise.resolve("first"),
  Promise.resolve("second"),
  Promise.resolve("third")
]).then(r => console.log("Race result:", r));

setTimeout(() => {
  console.log("\n=== Promise.any ===");
  Promise.any([
    Promise.reject("error 1"),
    Promise.resolve("WINNER"),
    Promise.resolve("too late")
  ])
    .then(r => console.log("Any result:", r))
    .catch(e => console.log("Any error:", e));
}, 100);
