// Very simple setTimeout test
console.log("Before setTimeout");

setTimeout(function() {
    console.log("Inside setTimeout callback!");
}, 100);

console.log("After setTimeout");
