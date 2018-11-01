const http = require('http')
const fs = require('fs')
const cluster = require('cluster')
const { webHookRun } = require('./webhook-runner.js')

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

  (new Promise(async (resolve, reject) => {
    try {
      await webHookRun(data)
      resolve()
    } catch(e) {
      reject(e)
    }
  })).catch((e) => {
    console.error(e)
  })
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
