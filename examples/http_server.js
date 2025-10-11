// http is available globally in Dougless Runtime - no require needed!

// Create a simple HTTP server
const server = http.createServer((req, res) => {
  console.log('Received request:', req.method, req.url)
  console.log('Headers:', req.headers)
  console.log('Body:', req.body)
  
  // Set response headers
  res.setHeader('Content-Type', 'application/json')
  res.setHeader('X-Custom-Header', 'Dougless-Runtime')
  
  // Send response
  res.statusCode = 200
  res.end(JSON.stringify({
    message: 'Hello from Dougless Runtime!',
    method: req.method,
    url: req.url,
    timestamp: Date.now()
  }))
})

// Start the server
server.listen(3000, () => {
  console.log('Server is running on http://localhost:3000')
  console.log('Test it with: curl http://localhost:3000')
  
  // Make a test GET request to our own server after a short delay
  setTimeout(() => {
    console.log('\n--- Making test GET request ---')
    http.get('http://localhost:3000/test', (err, response) => {
      if (err) {
        console.error('GET error:', err)
      } else {
        console.log('GET Status:', response.statusCode)
        console.log('GET Headers:', response.headers)
        console.log('GET Body:', response.body)
      }
    })
  }, 1000)
  
  // Make a test POST request after another delay
  setTimeout(() => {
    console.log('\n--- Making test POST request ---')
    http.post('http://localhost:3000/api/data', {
      name: 'Douglas',
      action: 'testing'
    }, (err, response) => {
      if (err) {
        console.error('POST error:', err)
      } else {
        console.log('POST Status:', response.statusCode)
        console.log('POST Headers:', response.headers)
        console.log('POST Body:', response.body)
      }
    })
  }, 2000)
})
