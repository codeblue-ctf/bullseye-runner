const http = require('http')
const { URL } = require('url')
const { calcScore } = require('./lib/calc-score.js')

process.on('message', async (data) => {
  const { id, callback_url, callback_authorization_token } = data
  const { succeeded, failed } = await calcScore(data)

  const postData = {
    id,
    succeeded,
    failed
  }
  const url = new URL(callback_url)
  const req = http.request({
    hostname: url.hostname,
    protocol: url.protocol,
    port: url.port,
    path: url.path,
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${callback_authorization_token}`
    }
  })
  req.write(JSON.stringify(postData))
  req.end()
})
