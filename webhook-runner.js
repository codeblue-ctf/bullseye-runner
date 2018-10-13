const http = require('http')
const { calcScore } = require('./lib/calc-score.js')

process.on('message', async (data) => {
  const { id, team, problem } = data
  const results = await calcScore(team, problem)

  const postData = {
    id,
    results
  }
  const req = http.request(process.env.BULLSEYE_WEB_WEBHOOK_URL, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    }
  })
  req.write(JSON.stringify(postData))
  req.end()
})
