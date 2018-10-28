const http = require('http')
const { calcScore } = require('./lib/calc-score.js')

process.on('message', async (data) => {
  const { id, callback_url, callback_authorization_token } = data
  const results = await calcScore(data)

  const postData = {
    id,
    results
  }
  const req = http.request(callback_url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${callback_authorization_token}`
    }
  })
  req.write(JSON.stringify(postData))
  req.end()
})
