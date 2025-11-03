// Simple HTTP server for benchmarking
const server = http.createServer((req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('Hello, World!');
});

server.listen(3456, () => {
    console.log('Benchmark server listening on port 3456');
});
