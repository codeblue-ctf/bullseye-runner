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
  const exploitContainer = `${config.registry.server}/${team.name}/${problem.exploit_container_name}`
  spawnSync('docker', ['pull', exploitContainer])
  spawnSync('docker', ['tag', exploitContainer, problem.exploit_container_name])

  // pull challenge container
  const challengeContainer = `${config.registry.server}/${problem.challenge_container_name}`
  spawnSync('docker', ['pull', challengeContainer])
  spawnSync('docker', ['tag', challengeContainer, problem.challenge_container_name])

  // pull flag submit container
  const flagSubmitContainer = `${config.registry.server}/${config.flagSubmitContainerName}`
  spawnSync('docker', ['pull', flagSubmitContainer])
  spawnSync('docker', ['tag', flagSubmitContainer, config.flagSubmitContainerName])
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
  const workingDir = fs.mkdtempSync('bulls-eye-runner')

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
