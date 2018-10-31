const http = require('http')
const { URL } = require('url')
const { calcScore } = require('./lib/calc-score.js')

process.on('message', async (data) => {
  const { id, callback_url, callback_authorization_token } = data
  const { succeeded, failed, imagesLog } = await calcScore(data)

  const postData = {
    schedule_uuid: id,
    succeeded,
    failed,
    error: imagesLog // XXX: return docker images result as error temporary
  }
  const url = new URL(callback_url)
  const req = http.request({
    hostname: url.hostname,
    protocol: url.protocol,
    path: url.pathname,
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${callback_authorization_token}`
    }
  })
  req.write(JSON.stringify(postData))
  req.end()
  console.debug('[postData]', JSON.stringify(postData))
})
