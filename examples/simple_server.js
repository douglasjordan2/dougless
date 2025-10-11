// http is available globally in Dougless Runtime - no require needed!

// Create a simple HTTP server
const server = http.createServer((req, res) => {
  console.log('Request received:', req.method, req.url)
  
  res.setHeader('Content-Type', 'application/json')
  res.setHeader('X-Powered-By', 'Dougless-Runtime')
  
  res.statusCode = 200
  res.end(JSON.stringify({
    message: 'Hello from Dougless Runtime!',
    method: req.method,
    url: req.url,
    timestamp: Date.now()
  }))
})

// Start the server and keep it running with a heartbeat
server.listen(3000, () => {
  console.log('🚀 Server running on http://localhost:3000')
  console.log('Press Ctrl+C to stop')
  
  // Keep the event loop alive with a repeating timer
  setInterval(() => {
    // This keeps the server running
  }, 10000)
})
