// WebSocket Server Example
// This example demonstrates creating a WebSocket server that echoes messages back to clients

console.log('Starting WebSocket Echo Server...');

// Create an HTTP server
const server = http.createServer(function(req, res) {
  // Serve a simple HTML page at the root
  if (req.url === '/') {
    res.setHeader('Content-Type', 'text/html');
    res.end(`
<!DOCTYPE html>
<html>
<head>
  <title>WebSocket Test</title>
  <style>
    body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
    #messages { border: 1px solid #ccc; height: 300px; overflow-y: scroll; padding: 10px; margin: 20px 0; }
    .message { margin: 5px 0; padding: 5px; }
    .sent { background: #e3f2fd; }
    .received { background: #f1f8e9; }
    .system { background: #fff3e0; color: #e65100; }
    input, button { padding: 10px; margin: 5px; }
    input { width: 300px; }
  </style>
</head>
<body>
  <h1>WebSocket Echo Test</h1>
  <div id="status">Connecting...</div>
  <div id="messages"></div>
  <input type="text" id="messageInput" placeholder="Type a message..." disabled />
  <button id="sendBtn" disabled>Send</button>
  <button id="closeBtn" disabled>Close Connection</button>

  <script>
    const messages = document.getElementById('messages');
    const status = document.getElementById('status');
    const input = document.getElementById('messageInput');
    const sendBtn = document.getElementById('sendBtn');
    const closeBtn = document.getElementById('closeBtn');

    // Connect to WebSocket
    const ws = new WebSocket('ws://localhost:8080/ws');

    function addMessage(text, type) {
      const div = document.createElement('div');
      div.className = 'message ' + type;
      div.textContent = text;
      messages.appendChild(div);
      messages.scrollTop = messages.scrollHeight;
    }

    ws.onopen = function() {
      status.textContent = 'Connected! (readyState: ' + ws.readyState + ')';
      status.style.color = 'green';
      input.disabled = false;
      sendBtn.disabled = false;
      closeBtn.disabled = false;
      addMessage('Connected to server', 'system');
    };

    ws.onmessage = function(event) {
      addMessage('Server: ' + event.data, 'received');
    };

    ws.onclose = function() {
      status.textContent = 'Disconnected (readyState: ' + ws.readyState + ')';
      status.style.color = 'red';
      input.disabled = true;
      sendBtn.disabled = true;
      closeBtn.disabled = true;
      addMessage('Disconnected from server', 'system');
    };

    ws.onerror = function(error) {
      addMessage('Error: ' + error, 'system');
      console.error('WebSocket error:', error);
    };

    sendBtn.onclick = function() {
      const message = input.value;
      if (message && ws.readyState === WebSocket.OPEN) {
        ws.send(message);
        addMessage('You: ' + message, 'sent');
        input.value = '';
      }
    };

    input.onkeypress = function(e) {
      if (e.key === 'Enter') {
        sendBtn.onclick();
      }
    };

    closeBtn.onclick = function() {
      ws.close();
    };
  </script>
</body>
</html>
    `);
  } else {
    res.statusCode = 404;
    res.end('Not Found');
  }
});

// Add WebSocket endpoint
server.websocket('/ws', {
  open: function(ws) {
    console.log('Client connected');
    console.log('  readyState:', ws.readyState, '(OPEN =', ws.OPEN + ')');
    
    // Send welcome message
    ws.send('Welcome to the echo server!');
  },

  message: function(msg) {
    console.log('Received message:', msg.data);
    console.log('  Message type:', msg.type);
    console.log('  Connection state:', msg); // Will show the message object
    
    // Echo the message back
    ws.send('Echo: ' + msg.data);
  },

  close: function() {
    console.log('Client disconnected');
  },

  error: function(err) {
    console.error('WebSocket error:', err);
  }
});

// Start the server
server.listen(8080, function() {
  console.log('Server is running on http://localhost:8080');
  console.log('WebSocket endpoint: ws://localhost:8080/ws');
  console.log('');
  console.log('Open http://localhost:8080 in your browser to test!');
  console.log('Press Ctrl+C to stop the server');
});
