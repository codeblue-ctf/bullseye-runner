const { spawnSync } = require('child_process')
const fs = require('fs')
const path = require('path')
const config = require('../config.js')
const { generateFlags } = require('./flag-generator.js')

const loginRegistry = () => {
  spawnSync('docker', [
    'login',
    config.registry.server,
    '-u',
    config.registry.admin.name,
    '-p',
    config.registry.admin.password
  ])
}

const pullImage = (team, problem) => {
  // pull exploit container
  const image = `${team.name}/${problem.exploit_container_name}`
  spawnSync('docker', ['pull', image])
  spawnSync('docker', ['tag', image, problem.exploit_container_name])

  // pull challenge container
  const image = problem.challenge_container_name
  spawnSync('docker', ['pull', image])

  // pull flag submit container
  const image = config.flagSubmitContainerName
  spawnSync('docker', ['pull', image])
}

const setFlag = (workingDir, flag) => {
  // XXX: flag filename is hardcoded
  const flagFile = path.join(workingDir, 'flag')
  fs.writeFileSync(flagFile, flag)
}

const runExploit = (workingDir, problem) => {
  spawnSync('docker-compose', ['up'], { cwd: workingDir })
  // TODO: sleep problem.exploit_timeout
  spawnSync('docker-compose', ['down'], { cwd: workingDir })
}

const getSubmittedFlags = (team, problem) => {
  // TODO: implement
}

const runExploits = (team, problem, flags) => {
  const workingDir = fs.mkdtempSync()

  loginRegistry()
  pullImage(team, problem)

  for (let i = 0; i < problem.calc_time; ++i) {
    const flag = flags[i]
    setFlag(workingDir, flag)
    runExploit(workingDir, problem)
  }
}

const calcScore = (team, problem) => {
  const flags = generateFlags(problem.exploit_trial_count)
  runExploits(team, problem, flags)
  const submittedFlags = getSubmittedFlags(team, problem)

  const correctFlags = flags.filter(flag => submittedFlags.includes(flag))
  return correctFlags.length
}

module.exports = {
  calcScore
}
