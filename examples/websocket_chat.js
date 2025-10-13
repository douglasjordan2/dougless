// WebSocket Chat Room Example
// This example demonstrates a multi-user chat room with broadcasting

console.log('Starting WebSocket Chat Server...');

// Store all connected clients
const clients = [];

const server = http.createServer(function(req, res) {
  if (req.url === '/') {
    res.setHeader('Content-Type', 'text/html');
    res.end(`
<!DOCTYPE html>
<html>
<head>
  <title>Chat Room</title>
  <style>
    body { 
      font-family: Arial, sans-serif; 
      max-width: 800px; 
      margin: 20px auto; 
      padding: 20px; 
      background: #f5f5f5;
    }
    h1 { color: #333; }
    #chat { 
      background: white; 
      border: 1px solid #ddd; 
      border-radius: 5px;
      height: 400px; 
      overflow-y: scroll; 
      padding: 15px; 
      margin: 20px 0;
      box-shadow: 0 2px 5px rgba(0,0,0,0.1);
    }
    .message { 
      margin: 8px 0; 
      padding: 8px 12px;
      border-radius: 4px;
      line-height: 1.4;
    }
    .user-msg { 
      background: #e3f2fd; 
      border-left: 3px solid #2196f3;
    }
    .own-msg { 
      background: #f1f8e9; 
      border-left: 3px solid #4caf50;
    }
    .system { 
      background: #fff3e0; 
      color: #e65100;
      font-style: italic;
      text-align: center;
      border-left: none;
    }
    .input-area {
      display: flex;
      gap: 10px;
    }
    #username {
      width: 150px;
      padding: 10px;
      border: 1px solid #ddd;
      border-radius: 4px;
    }
    #messageInput { 
      flex: 1;
      padding: 10px;
      border: 1px solid #ddd;
      border-radius: 4px;
    }
    button { 
      padding: 10px 20px;
      background: #2196f3;
      color: white;
      border: none;
      border-radius: 4px;
      cursor: pointer;
    }
    button:hover { background: #1976d2; }
    button:disabled { 
      background: #ccc; 
      cursor: not-allowed;
    }
    #status {
      padding: 10px;
      margin-bottom: 10px;
      border-radius: 4px;
      text-align: center;
      font-weight: bold;
    }
    .connected { background: #c8e6c9; color: #2e7d32; }
    .disconnected { background: #ffcdd2; color: #c62828; }
  </style>
</head>
<body>
  <h1>💬 WebSocket Chat Room</h1>
  <div id="status" class="disconnected">Connecting...</div>
  <div id="chat"></div>
  <div class="input-area">
    <input type="text" id="username" placeholder="Your name..." />
    <input type="text" id="messageInput" placeholder="Type a message..." disabled />
    <button id="sendBtn" disabled>Send</button>
  </div>

  <script>
    const chat = document.getElementById('chat');
    const status = document.getElementById('status');
    const usernameInput = document.getElementById('username');
    const messageInput = document.getElementById('messageInput');
    const sendBtn = document.getElementById('sendBtn');

    const ws = new WebSocket('ws://localhost:8080/chat');

    function addMessage(text, className) {
      const div = document.createElement('div');
      div.className = 'message ' + className;
      div.textContent = text;
      chat.appendChild(div);
      chat.scrollTop = chat.scrollHeight;
    }

    ws.onopen = function() {
      status.textContent = '✓ Connected to chat';
      status.className = 'connected';
      messageInput.disabled = false;
      sendBtn.disabled = false;
      addMessage('Connected to chat room', 'system');
    };

    ws.onmessage = function(event) {
      const data = JSON.parse(event.data);
      
      if (data.type === 'system') {
        addMessage(data.message, 'system');
      } else if (data.type === 'chat') {
        const isOwn = data.username === usernameInput.value;
        const className = isOwn ? 'own-msg' : 'user-msg';
        addMessage(data.username + ': ' + data.message, className);
      }
    };

    ws.onclose = function() {
      status.textContent = '✗ Disconnected';
      status.className = 'disconnected';
      messageInput.disabled = true;
      sendBtn.disabled = true;
      addMessage('Disconnected from chat room', 'system');
    };

    ws.onerror = function(error) {
      addMessage('Error: ' + error, 'system');
    };

    function sendMessage() {
      const username = usernameInput.value.trim() || 'Anonymous';
      const message = messageInput.value.trim();
      
      if (message && ws.readyState === WebSocket.OPEN) {
        const data = JSON.stringify({
          username: username,
          message: message
        });
        ws.send(data);
        messageInput.value = '';
      }
    }

    sendBtn.onclick = sendMessage;
    messageInput.onkeypress = function(e) {
      if (e.key === 'Enter') sendMessage();
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

// WebSocket chat endpoint
server.websocket('/chat', {
  open: function(ws) {
    console.log('New client connected. Total clients:', clients.length + 1);
    
    // Add client to list
    clients.push(ws);
    
    // Notify everyone
    const joinMsg = JSON.stringify({
      type: 'system',
      message: 'A new user joined the chat (' + clients.length + ' users online)'
    });
    
    clients.forEach(function(client) {
      if (client.readyState === ws.OPEN) {
        client.send(joinMsg);
      }
    });
  },

  message: function(msg) {
    const data = JSON.parse(msg.data);
    console.log('Message from', data.username + ':', data.message);
    
    // Broadcast to all clients
    const broadcastMsg = JSON.stringify({
      type: 'chat',
      username: data.username,
      message: data.message
    });
    
    clients.forEach(function(client) {
      if (client.readyState === client.OPEN) {
        client.send(broadcastMsg);
      }
    });
  },

  close: function() {
    // Find and remove disconnected client
    for (let i = 0; i < clients.length; i++) {
      if (clients[i].readyState === clients[i].CLOSED || 
          clients[i].readyState === clients[i].CLOSING) {
        clients.splice(i, 1);
        break;
      }
    }
    
    console.log('Client disconnected. Remaining clients:', clients.length);
    
    // Notify remaining clients
    const leaveMsg = JSON.stringify({
      type: 'system',
      message: 'A user left the chat (' + clients.length + ' users online)'
    });
    
    clients.forEach(function(client) {
      if (client.readyState === client.OPEN) {
        client.send(leaveMsg);
      }
    });
  },

  error: function(err) {
    console.error('WebSocket error:', err);
  }
});

server.listen(8080, function() {
  console.log('Chat server is running on http://localhost:8080');
  console.log('Open multiple browser tabs to test multi-user chat!');
  console.log('Press Ctrl+C to stop the server');
});
