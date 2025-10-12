# HTTP API Guide

## Overview

Dougless provides a unique global `http` API for HTTP client and server operations. Unlike Node.js which requires `require('http')`, the `http` object is always available globally.

## Why Global?

**Dougless Philosophy**: HTTP operations are fundamental to modern web applications and should be as accessible as `console`. This makes Dougless code cleaner and more intuitive for web development.

**Comparison:**

```javascript
// Node.js
const http = require('http');
http.get(url, callback);

// Dougless
http.get(url, callback);  // No require!
```

---

## HTTP Client Operations

### `http.get(url, callback)`

Make an HTTP GET request asynchronously.

**Parameters:**
- `url` (string) - Full URL to request
- `callback` (function) - Callback function `(err, response)`

**Response Object:**
- `statusCode` (number) - HTTP status code (e.g., 200, 404)
- `status` (string) - HTTP status text (e.g., "200 OK")
- `headers` (object) - Response headers
- `body` (string) - Response body as a string

**Example:**
```javascript
http.get('https://api.example.com/data', function(err, response) {
    if (err) {
        console.error('Request failed:', err);
        return;
    }
    
    console.log('Status:', response.statusCode);
    console.log('Headers:', response.headers);
    console.log('Body:', response.body);
    
    // Parse JSON if needed
    const data = JSON.parse(response.body);
    console.log('Data:', data);
});
```

---

### `http.post(url, data, callback)`

Make an HTTP POST request with JSON data asynchronously.

**Parameters:**
- `url` (string) - Full URL to request
- `data` (object) - JavaScript object to send as JSON
- `callback` (function) - Callback function `(err, response)`

**Content Type:**
- Default: `application/json`
- Data is automatically JSON-encoded

**Response Object:**
Same as `http.get()` - contains `statusCode`, `status`, `headers`, and `body`.

**Example:**
```javascript
const payload = {
    username: 'dougless',
    action: 'create',
    data: {
        name: 'My Project',
        type: 'web-app'
    }
};

http.post('https://api.example.com/projects', payload, function(err, response) {
    if (err) {
        console.error('POST failed:', err);
        return;
    }
    
    console.log('Created! Status:', response.statusCode);
    console.log('Response:', response.body);
});
```

---

## HTTP Server Operations

### `http.createServer(requestHandler)`

Create an HTTP server.

**Parameters:**
- `requestHandler` (function) - Function called for each request: `(req, res)`

**Returns:**
- Server object with `listen()` method

**Request Object (req):**
- `method` (string) - HTTP method (GET, POST, PUT, DELETE, etc.)
- `url` (string) - Request URL path
- `headers` (object) - Request headers
- `body` (string) - Request body content

**Response Object (res):**
- `statusCode` (number) - Set the status code (default: 200)
- `setHeader(name, value)` - Set a response header
- `end(data)` - Send response and close connection

**Example:**
```javascript
const server = http.createServer(function(req, res) {
    console.log('Request:', req.method, req.url);
    
    // Set response headers
    res.setHeader('Content-Type', 'application/json');
    res.setHeader('X-Powered-By', 'Dougless-Runtime');
    
    // Set status code
    res.statusCode = 200;
    
    // Send response
    res.end(JSON.stringify({
        message: 'Hello from Dougless!',
        method: req.method,
        path: req.url
    }));
});
```

---

### `server.listen(port, [callback])`

Start the HTTP server listening on a port.

**Parameters:**
- `port` (string|number) - Port number to listen on
- `callback` (function) - Optional callback when server starts

**Example:**
```javascript
server.listen(3000, function() {
    console.log('Server running on http://localhost:3000');
});
```

**Note:** The server runs in a background goroutine. Keep the event loop alive with timers if needed.

---

## Complete Examples

### Example 1: Simple GET Request

```javascript
console.log('Fetching data...');

http.get('https://jsonplaceholder.typicode.com/todos/1', function(err, response) {
    if (err) {
        console.error('Error:', err);
        return;
    }
    
    console.log('Status:', response.statusCode);
    
    const todo = JSON.parse(response.body);
    console.log('Todo:', todo.title);
    console.log('Completed:', todo.completed);
});
```

### Example 2: POST Request with Data

```javascript
const userData = {
    name: 'Douglas Jordan',
    email: 'doug@example.com',
    role: 'developer'
};

http.post('https://api.example.com/users', userData, function(err, response) {
    if (err) {
        console.error('Failed to create user:', err);
        return;
    }
    
    if (response.statusCode === 201) {
        console.log('User created successfully!');
        const newUser = JSON.parse(response.body);
        console.log('New user ID:', newUser.id);
    } else {
        console.error('Unexpected status:', response.statusCode);
    }
});
```

### Example 3: Basic HTTP Server

```javascript
const server = http.createServer(function(req, res) {
    console.log('Request received:', req.method, req.url);
    
    // Set headers
    res.setHeader('Content-Type', 'application/json');
    res.setHeader('Access-Control-Allow-Origin', '*');
    
    // Handle different routes
    if (req.url === '/health') {
        res.statusCode = 200;
        res.end(JSON.stringify({ status: 'healthy' }));
    } else if (req.url === '/api/data') {
        res.statusCode = 200;
        res.end(JSON.stringify({
            message: 'API endpoint',
            timestamp: Date.now()
        }));
    } else {
        res.statusCode = 404;
        res.end(JSON.stringify({ error: 'Not found' }));
    }
});

server.listen(3000, function() {
    console.log('Server listening on http://localhost:3000');
});

// Keep server running
setInterval(function() {}, 10000);
```

### Example 4: Server with Request Body Handling

```javascript
const server = http.createServer(function(req, res) {
    res.setHeader('Content-Type', 'application/json');
    
    if (req.method === 'POST' && req.url === '/api/echo') {
        // Echo back the request body
        console.log('Received body:', req.body);
        
        let requestData;
        try {
            requestData = JSON.parse(req.body);
        } catch (e) {
            res.statusCode = 400;
            res.end(JSON.stringify({ error: 'Invalid JSON' }));
            return;
        }
        
        res.statusCode = 200;
        res.end(JSON.stringify({
            echo: requestData,
            receivedAt: Date.now()
        }));
    } else {
        res.statusCode = 405;
        res.end(JSON.stringify({ error: 'Method not allowed' }));
    }
});

server.listen(8080, function() {
    console.log('Echo server running on http://localhost:8080');
});
```

### Example 5: Full-Stack Example (Server + Client)

```javascript
// Create server
const server = http.createServer(function(req, res) {
    res.setHeader('Content-Type', 'application/json');
    
    if (req.url === '/api/message') {
        res.statusCode = 200;
        res.end(JSON.stringify({
            message: 'Hello from the server!',
            time: new Date().toISOString()
        }));
    } else {
        res.statusCode = 404;
        res.end(JSON.stringify({ error: 'Not found' }));
    }
});

// Start server
server.listen(3000, function() {
    console.log('Server started on http://localhost:3000');
    
    // Make a request to our own server after a delay
    setTimeout(function() {
        console.log('\nMaking request to our server...');
        
        http.get('http://localhost:3000/api/message', function(err, response) {
            if (err) {
                console.error('Request failed:', err);
            } else {
                console.log('Response:', response.body);
            }
        });
    }, 1000);
});

// Keep running
setInterval(function() {}, 10000);
```

---

## Error Handling

HTTP client operations follow the Node.js error-first callback pattern:

```javascript
function callback(err, response) {
    if (err) {
        // Network error, timeout, or other failure
        console.error('Error:', err);
        return;
    }
    
    // Check HTTP status code
    if (response.statusCode >= 400) {
        console.error('HTTP error:', response.statusCode, response.status);
        return;
    }
    
    // Success - use response
    console.log('Success:', response.body);
}
```

**Common Errors:**
- Network unreachable
- Connection timeout
- DNS resolution failure
- Invalid URL
- Server not responding

**HTTP Status Codes:**
- `2xx` - Success (200 OK, 201 Created, etc.)
- `3xx` - Redirect (301, 302, etc.)
- `4xx` - Client error (404 Not Found, 400 Bad Request, etc.)
- `5xx` - Server error (500 Internal Server Error, 503 Service Unavailable, etc.)

---

## Working with Headers

### Reading Response Headers

```javascript
http.get('https://api.example.com/data', function(err, response) {
    if (!err) {
        console.log('Content-Type:', response.headers['Content-Type']);
        console.log('All headers:', response.headers);
    }
});
```

### Setting Request Headers

Currently, `http.get()` and `http.post()` use default headers. Custom headers support is planned for future versions.

### Setting Response Headers (Server)

```javascript
const server = http.createServer(function(req, res) {
    // Set single header
    res.setHeader('Content-Type', 'application/json');
    
    // Set multiple headers
    res.setHeader('X-Custom-Header', 'Custom Value');
    res.setHeader('Access-Control-Allow-Origin', '*');
    res.setHeader('Cache-Control', 'no-cache');
    
    res.end('Response data');
});
```

---

## Server Patterns

### Keep-Alive Pattern

Servers need to keep the event loop alive:

```javascript
const server = http.createServer(requestHandler);

server.listen(3000, function() {
    console.log('Server started');
    
    // Keep event loop alive
    setInterval(function() {
        // Heartbeat - prevents server from exiting
    }, 10000);
});
```

### Routing Pattern

```javascript
const server = http.createServer(function(req, res) {
    res.setHeader('Content-Type', 'application/json');
    
    // Simple routing
    if (req.url === '/') {
        res.statusCode = 200;
        res.end(JSON.stringify({ message: 'Home' }));
    } else if (req.url === '/api/users') {
        res.statusCode = 200;
        res.end(JSON.stringify({ users: [] }));
    } else if (req.url === '/api/posts') {
        res.statusCode = 200;
        res.end(JSON.stringify({ posts: [] }));
    } else {
        res.statusCode = 404;
        res.end(JSON.stringify({ error: 'Not found' }));
    }
});
```

### Method-Based Routing

```javascript
const server = http.createServer(function(req, res) {
    res.setHeader('Content-Type', 'application/json');
    
    if (req.url === '/api/data') {
        if (req.method === 'GET') {
            // Handle GET
            res.statusCode = 200;
            res.end(JSON.stringify({ data: 'Here it is' }));
        } else if (req.method === 'POST') {
            // Handle POST
            const body = JSON.parse(req.body);
            res.statusCode = 201;
            res.end(JSON.stringify({ created: body }));
        } else if (req.method === 'DELETE') {
            // Handle DELETE
            res.statusCode = 204;
            res.end();
        } else {
            res.statusCode = 405;
            res.end(JSON.stringify({ error: 'Method not allowed' }));
        }
    } else {
        res.statusCode = 404;
        res.end(JSON.stringify({ error: 'Not found' }));
    }
});
```

---

## Limitations

### Current Limitations (Phase 3)
- GET and POST only (no PUT, DELETE, PATCH methods yet)
- No custom request headers
- No streaming (entire body loaded into memory)
- No HTTPS certificate validation options
- No request timeout configuration
- Server cannot be stopped programmatically

### Future Enhancements (Phase 4+)
- Full HTTP method support (PUT, DELETE, PATCH, OPTIONS)
- Custom request headers
- Streaming request/response bodies
- Request timeout configuration
- Connection pooling
- HTTPS server support
- Server shutdown/restart capabilities
- Middleware support
- WebSocket upgrade (Phase 4)

---

## Comparison with Node.js

| Feature | Node.js | Dougless |
|---------|---------|----------|
| **Import** | `require('http')` | Global (no require) |
| **GET** | `http.get()` | `http.get()` |
| **POST** | `http.request()` + config | `http.post()` (simpler) |
| **Server** | `http.createServer()` | `http.createServer()` |
| **Listen** | `server.listen()` | `server.listen()` |
| **Request object** | Full Node.js req | Simplified (method, url, headers, body) |
| **Response object** | Full Node.js res | Simplified (statusCode, setHeader, end) |
| **Streaming** | ✅ Yes | ❌ Not yet |
| **All HTTP methods** | ✅ Yes | ⏳ GET/POST only |

---

## Best Practices

1. **Always handle errors** - Check the error parameter in callbacks
2. **Validate response status** - Don't assume 2xx, check `statusCode`
3. **Parse JSON safely** - Use try/catch when parsing response bodies
4. **Set appropriate headers** - Always set Content-Type for responses
5. **Keep servers alive** - Use `setInterval()` to prevent server exit
6. **Handle all routes** - Return 404 for unknown paths
7. **Log requests** - Console.log incoming requests for debugging
8. **Validate request bodies** - Check and parse POST data carefully

---

## Testing Your HTTP Code

### Test Client Requests
```bash
# From another terminal or external tool
curl http://localhost:3000
curl -X POST http://localhost:3000/api/data -d '{"test":"data"}'
```

### Test External APIs
```javascript
// Use a public API for testing
http.get('https://jsonplaceholder.typicode.com/posts/1', function(err, res) {
    if (!err) console.log('Success!', res.body);
});
```

---

**Happy HTTP coding with Dougless!** 🚀
