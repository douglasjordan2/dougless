// =======================================
// Dougless Runtime - HTTP Features Demo
// =======================================
// Demonstrates the global http API (no require needed!)
// Includes both HTTP client and server functionality

console.log('=== HTTP Features Demo ===\n');
console.log('Starting HTTP server and client demonstrations...\n');

// ======================================
// PART 1: HTTP SERVER
// ======================================

console.log('PART 1: Creating HTTP Server');
console.log('----------------------------');

const server = http.createServer(function(req, res) {
    console.log('  → Request received:', req.method, req.url);
    
    // Set common headers
    res.setHeader('Content-Type', 'application/json');
    res.setHeader('X-Powered-By', 'Dougless-Runtime');
    
    // Route handling
    if (req.url === '/') {
        // Home route
        res.statusCode = 200;
        res.end(JSON.stringify({
            message: 'Welcome to Dougless HTTP Server!',
            endpoints: ['/api/data', '/api/echo', '/health'],
            timestamp: Date.now()
        }));
        
    } else if (req.url === '/health') {
        // Health check
        res.statusCode = 200;
        res.end(JSON.stringify({
            status: 'healthy',
            uptime: Date.now()
        }));
        
    } else if (req.url === '/api/data') {
        // API data endpoint
        res.statusCode = 200;
        res.end(JSON.stringify({
            data: [
                { id: 1, name: 'Item 1' },
                { id: 2, name: 'Item 2' },
                { id: 3, name: 'Item 3' }
            ],
            count: 3
        }));
        
    } else if (req.url === '/api/echo' && req.method === 'POST') {
        // Echo endpoint - returns request body
        console.log('  → Request body:', req.body);
        
        let requestData;
        try {
            requestData = JSON.parse(req.body);
        } catch (e) {
            res.statusCode = 400;
            res.end(JSON.stringify({
                error: 'Invalid JSON',
                message: e.message
            }));
            return;
        }
        
        res.statusCode = 200;
        res.end(JSON.stringify({
            echo: requestData,
            receivedAt: Date.now()
        }));
        
    } else {
        // 404 Not Found
        res.statusCode = 404;
        res.end(JSON.stringify({
            error: 'Not Found',
            path: req.url
        }));
    }
});

// Start the server
server.listen(3000, function() {
    console.log('✓ Server started on http://localhost:3000');
    console.log('✓ Server is ready to accept requests\n');
    
    // ======================================
    // PART 2: HTTP CLIENT - GET REQUESTS
    // ======================================
    
    setTimeout(function() {
        console.log('\nPART 2: HTTP Client - GET Requests');
        console.log('-----------------------------------');
        
        // Test 1: GET home route
        console.log('\n1. GET http://localhost:3000/');
        http.get('http://localhost:3000/', function(err, response) {
            if (err) {
                console.error('   ✗ Error:', err);
            } else {
                console.log('   ✓ Status:', response.statusCode);
                console.log('   ✓ Response:', response.body);
            }
        });
    }, 500);
    
    setTimeout(function() {
        // Test 2: GET health check
        console.log('\n2. GET http://localhost:3000/health');
        http.get('http://localhost:3000/health', function(err, response) {
            if (err) {
                console.error('   ✗ Error:', err);
            } else {
                console.log('   ✓ Status:', response.statusCode);
                console.log('   ✓ Response:', response.body);
            }
        });
    }, 1000);
    
    setTimeout(function() {
        // Test 3: GET API data
        console.log('\n3. GET http://localhost:3000/api/data');
        http.get('http://localhost:3000/api/data', function(err, response) {
            if (err) {
                console.error('   ✗ Error:', err);
            } else {
                console.log('   ✓ Status:', response.statusCode);
                const data = JSON.parse(response.body);
                console.log('   ✓ Received', data.count, 'items');
                console.log('   ✓ First item:', data.data[0].name);
            }
        });
    }, 1500);
    
    setTimeout(function() {
        // Test 4: GET non-existent route (404)
        console.log('\n4. GET http://localhost:3000/nonexistent (expect 404)');
        http.get('http://localhost:3000/nonexistent', function(err, response) {
            if (err) {
                console.error('   ✗ Error:', err);
            } else {
                console.log('   ✓ Status:', response.statusCode, '(404 as expected)');
                console.log('   ✓ Response:', response.body);
            }
        });
    }, 2000);
    
    // ======================================
    // PART 3: HTTP CLIENT - POST REQUESTS
    // ======================================
    
    setTimeout(function() {
        console.log('\n\nPART 3: HTTP Client - POST Requests');
        console.log('------------------------------------');
        
        // Test 5: POST with JSON data
        console.log('\n5. POST http://localhost:3000/api/echo');
        const postData = {
            name: 'Douglas',
            action: 'testing',
            items: [1, 2, 3],
            metadata: {
                timestamp: Date.now(),
                source: 'dougless-runtime'
            }
        };
        
        console.log('   → Sending:', JSON.stringify(postData));
        http.post('http://localhost:3000/api/echo', postData, function(err, response) {
            if (err) {
                console.error('   ✗ Error:', err);
            } else {
                console.log('   ✓ Status:', response.statusCode);
                console.log('   ✓ Response:', response.body);
            }
        });
    }, 2500);
    
    setTimeout(function() {
        // Test 6: POST to non-POST endpoint
        console.log('\n6. POST http://localhost:3000/ (GET-only endpoint)');
        http.post('http://localhost:3000/', { test: 'data' }, function(err, response) {
            if (err) {
                console.error('   ✗ Error:', err);
            } else {
                console.log('   ✓ Status:', response.statusCode);
                console.log('   ✓ Server handled it gracefully');
            }
        });
    }, 3000);
    
    // ======================================
    // PART 4: EXTERNAL API TEST
    // ======================================
    
    setTimeout(function() {
        console.log('\n\nPART 4: External API Request');
        console.log('-----------------------------');
        console.log('\n7. GET https://jsonplaceholder.typicode.com/posts/1');
        
        http.get('https://jsonplaceholder.typicode.com/posts/1', function(err, response) {
            if (err) {
                console.error('   ✗ Error:', err);
            } else {
                console.log('   ✓ Status:', response.statusCode);
                const post = JSON.parse(response.body);
                console.log('   ✓ Post title:', post.title);
                console.log('   ✓ User ID:', post.userId);
            }
        });
    }, 3500);
    
    // ======================================
    // FINAL MESSAGE
    // ======================================
    
    setTimeout(function() {
        console.log('\n\n=== HTTP Demo Complete ===');
        console.log('All HTTP features demonstrated successfully!');
        console.log('\nServer is still running on http://localhost:3000');
        console.log('You can test it with: curl http://localhost:3000');
        console.log('\nPress Ctrl+C to stop the server');
    }, 4500);
    
    // Keep the server running
    setInterval(function() {
        // Heartbeat to keep event loop alive
    }, 10000);
});
