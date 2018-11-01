const http = require('http')
const fs = require('fs')

const validData = (data) => {
  params = [
    'id',
    'trials_count',
    'timeout',
    'docker_compose',
    'registry_host',
    'admin_username',
    'admin_password',
    'flag_template',
    'callback_url',
    'callback_authorization_token'
  ]
  const missingparam = params.find((el) => {
    data[el] === undefined
  })

  return missingparam === undefined
}

const processData = (data) => {
  if (!validData(data)) throw 'invalid data';

  const child = require('child_process').fork('./webhook-runner')
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
      console.error(e)
      res.end(JSON.stringify({'result': 'failed'}))
    }
  })
})
server.listen(3000)
