// Simple test script for dougless runtime
console.log("Hello from Dougless Runtime!");
console.warn("This is a warning");
console.error("This is an error");

// Test require (placeholder modules)
try {
    const fs = require('fs');
    const http = require('http');
    const path = require('path');
    console.log("Modules loaded successfully");
} catch (e) {
    console.error("Failed to load modules:", e.message);
}
