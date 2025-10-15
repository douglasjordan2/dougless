console.log("Starting...");

const p1 = Promise.resolve("first");
const p2 = Promise.resolve("second");
const p3 = Promise.resolve("third");

console.log("Attaching handlers...");
p1.then(r => console.log("p1 handler:", r));
p2.then(r => console.log("p2 handler:", r));
p3.then(r => console.log("p3 handler:", r));

console.log("Handlers attached");

setTimeout(() => {
  console.log("\nUsing Promise.race:");
  Promise.race([p1, p2, p3]).then(r => console.log("Race result:", r));
}, 100);
