console.log('Testing HTTP GET...')

http.get('https://jsonplaceholder.typicode.com/todos/1', (err, response) => {
  if (err) {
    console.error('Error:', err)
  } else {
    console.log('Status:', response.statusCode)
    console.log('Body:', response.body)
  }
})
