const { spawnSync } = require('child_process')
const fs = require('fs')
const path = require('path')
const config = require('../config.js')
const { generateFlags } = require('./flag-generator.js')

const loginRegistry = () => {
  const { stdout, stderr, status } = spawnSync('docker', [
    'login',
    config.registry.server,
    '-u',
    config.registry.admin.name,
    '-p',
    config.registry.admin.password
  ])
  console.debug('stdout', stdout.toString())
  console.debug('stderr', stderr.toString())
}

const pullImage = (team, problem) => {
  // pull exploit container
  const exploitContainer = `${config.registry.server}/${team.name}/${problem.exploit_container_name}`
  spawnSync('docker', ['pull', exploitContainer])
  spawnSync('docker', ['tag', exploitContainer, problem.exploit_container_name])
  console.debug('docker pull', exploitContainer)

  // pull challenge container
  const challengeContainer = `${config.registry.server}/${problem.problem_container_name}`
  spawnSync('docker', ['pull', challengeContainer])
  spawnSync('docker', ['tag', challengeContainer, problem.problem_container_name])
  console.debug('docker pull', challengeContainer)

  // pull flag submit container
  const flagSubmitContainer = `${config.registry.server}/${config.flagSubmitContainerName}`
  spawnSync('docker', ['pull', flagSubmitContainer])
  spawnSync('docker', ['tag', flagSubmitContainer, config.flagSubmitContainerName])
  console.debug('docker pull', flagSubmitContainer)
}

const setFlag = (workingDir, flag) => {
  // XXX: flag filename is hardcoded
  const flagFile = path.join(workingDir, 'flag')
  fs.writeFileSync(flagFile, flag)
}

const runExploit = (i, workingDir, problem) => {
  console.debug('start exploit', i)
  spawnSync('docker-compose', ['up'], { cwd: workingDir })
  // TODO: sleep problem.exploit_timeout
  spawnSync('docker-compose', ['down'], { cwd: workingDir })
  console.debug('end exploit', i)
}

const getSubmittedFlags = (team, problem) => {
  // TODO: implement
}

const runExploits = (team, problem, flags) => {
  const workingDir = fs.mkdtempSync('bulls-eye-runner')
  console.debug('workingDir', workingDir)

  loginRegistry()
  pullImage(team, problem)

  for (let i = 0; i < problem.exploit_trial_count; ++i) {
    const flag = flags[i]
    setFlag(workingDir, flag)
    runExploit(i, workingDir, problem)
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
