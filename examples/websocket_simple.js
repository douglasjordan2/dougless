// Minimal WebSocket Server Example
// The simplest possible WebSocket server

console.log('Starting minimal WebSocket server...');

const server = http.createServer(function(req, res) {
  res.end('WebSocket server is running. Connect to ws://localhost:8080/ws');
});

server.websocket('/ws', {
  open: function(ws) {
    console.log('✓ Client connected');
    ws.send('Hello from Dougless!');
  },
  
  message: function(msg) {
    console.log('← Received:', msg.data);
    ws.send('You said: ' + msg.data);
  },
  
  close: function() {
    console.log('✗ Client disconnected');
  },
  
  error: function(err) {
    console.error('Error:', err);
  }
});

server.listen(8080, function() {
  console.log('Server running on http://localhost:8080');
  console.log('WebSocket at ws://localhost:8080/ws');
  console.log('Press Ctrl+C to stop');
});
