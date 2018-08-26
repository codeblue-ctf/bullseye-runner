const { spawnSync } = require('child_proces')
const config = require('../config.js')
const { generateFlags } = require('./lib/flagGenerator.js')

const loginToRegistry = () => {
  spawnSync('docker', [
    'login',
    config.registry.server,
    '-u',
    config.registry.usre.name,
    '-p',
    config.registry.usre.password
  ])
}

const pullImage = (teamName, problemName) => {
  const image = `${teamName}/${problemName}`
  spawnSync('docker', ['pull', image])
}

const setFlag = (teamName, problemName, flag) => {
  // TODO: implement
}

const runExploit = (teamName, problemName) => {
  // TODO: implement
}

const getSubmittedFlags = (teamName, problemName) => {
  // TODO: implement
}

const runExploits = (teamName, problemName, calcTime, runTime, flags) => {
  loginToRegistry()
  // TODO: set tag
  pullImage(teamName, problemName)

  for (let i = 0; i < calcTime; ++i) {
    const flag = flags[i]
    setFlag(teamName, problemName, flag)
    runExploit(teamName, problemName)
  }
}

const calcScore = (teamName, problemName, calcTime, runTime) => {
  const flags = generateFlags(calcTime)
  runExploits(teamName, problemName, calcTime, runTime)
  const submittedFlags = getSubmittedFlags(teamName, problemName)

  const correctFlags = flags.filter(flag => submittedFlags.includes(flag))
  return correctFlags.length
}
