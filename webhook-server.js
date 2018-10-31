const http = require('http')
const fs = require('fs')
const child_process = require('child_process')

const processData = (data) => {
  const child = child_process.fork('./webhook-runner')
  child.send(data)
}

const server = http.createServer((req, res) => {
  res.writeHead(200, { 'Content-Type': 'application/json' })
  let rawData = ''
  req.on('data', (chunk) => {
    rawData += chunk;
  })
  req.on('end', () => {
    try {
      console.debug('[receivedData]', rawData)
      const data = JSON.parse(rawData)
      processData(data)

      res.end(JSON.stringify({'result': 'success'}))
    } catch(e) {
      res.end(JSON.stringify({'result': 'failed'}))
    }
  })
})
server.listen(3000)
