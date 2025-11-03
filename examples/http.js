// Test HTTP operations
console.log('Testing HTTP...');

// Test HTTP GET
console.log('Making GET request...');
const getResponse = http.get('https://jsonplaceholder.typicode.com/posts/1');

setTimeout(() => {
  console.log('GET Status:', getResponse.statusCode);
  console.log('GET Body:', getResponse.body);
}, 1000);

// Test HTTP POST
console.log('Making POST request...');
const postResponse = http.post('https://jsonplaceholder.typicode.com/posts', {
  title: 'Test Post',
  body: 'This is a test',
  userId: 1
});

setTimeout(() => {
  console.log('POST Status:', postResponse.statusCode);
  console.log('POST Body:', postResponse.body);
}, 1500);

// Test HTTP server
console.log('Creating HTTP server...');
const server = http.createServer((req, res) => {
  console.log(`Server received ${req.method} ${req.url}`);
  
  res.writeHead(200, { 'Content-Type': 'text/plain' });
  res.end('Hello from Dougless HTTP server!');
});

server.listen(8080, () => {
  console.log('Server listening on http://localhost:8080');
  
  // Test the server with a GET request
  setTimeout(() => {
    console.log('Testing server with GET request...');
    const testResponse = http.get('http://localhost:8080/test');
    
    setTimeout(() => {
      console.log('Server test status:', testResponse.statusCode);
      console.log('Server test body:', testResponse.body);
      
      // Close the server
      server.close();
      console.log('Server closed');
    }, 500);
  }, 100);
});

console.log('HTTP operations started...');
