const { spawnSync } = require('child_process')
const fs = require('fs')
const path = require('path')
const config = require('../config.js')
const { generateFlags } = require('./flag-generator.js')
const { sleep } = require('./sleep.js')

const spawnSyncWithLog = (command, args, options) => {
  console.debug('[spawnSync]', command, args)
  const result = spawnSync(command, args, options)
  console.debug('[spawnSync][stdout]', result.stdout.toString())
  console.debug('[spawnSync][stderr]', result.stderr.toString())
  return result
}

const loginRegistry = () => {
  spawnSyncWithLog('docker', [
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
  spawnSyncWithLog('docker', ['pull', exploitContainer])
  spawnSyncWithLog('docker', ['tag', exploitContainer, problem.exploit_container_name])

  // pull challenge container
  const challengeContainer = `${config.registry.server}/${problem.problem_container_name}`
  spawnSyncWithLog('docker', ['pull', challengeContainer])
  spawnSyncWithLog('docker', ['tag', challengeContainer, problem.problem_container_name])

  // pull flag submit container
  const flagSubmitContainer = `${config.registry.server}/${config.flagSubmitContainerName}`
  spawnSyncWithLog('docker', ['pull', flagSubmitContainer])
  spawnSyncWithLog('docker', ['tag', flagSubmitContainer, config.flagSubmitContainerName])
}

const setupDockerCompose = (workingDir, problem) => {
  const dockerComposeFile = path.join(workingDir, 'docker-compose.yml')
  fs.writeFileSync(dockerComposeFile, problem.docker_compose)

  spawnSyncWithLog('docker-compose', ['build'], { cwd: workingDir })
}

const setFlag = (workingDir, flag) => {
  // XXX: flag filename is hardcoded
  const flagFile = path.join(workingDir, 'flag')
  fs.writeFileSync(flagFile, flag)
}

const runExploit = async (i, workingDir, problem) => {
  console.debug('start exploit', i)

  // XXX: It should be `spawn` instead of `spawnSync`?
  // Malicious exploit container may take times to start up
  spawnSyncWithLog('docker-compose', ['up', '-d'], { cwd: workingDir })

  await sleep(problem.exploit_timeout * 1000)

  spawnSyncWithLog('docker-compose', ['down'], { cwd: workingDir })

  console.debug('end exploit', i)
}

const getSubmittedFlags = (team, problem) => {
  // TODO: implement
}

const runExploits = async (team, problem, flags) => {
  const workingDir = fs.mkdtempSync('bulls-eye-runner')
  console.debug('workingDir', workingDir)

  loginRegistry()
  pullImage(team, problem)
  setupDockerCompose(workingDir, problem)

  for (let i = 0; i < problem.exploit_trial_count; ++i) {
    const flag = flags[i]
    setFlag(workingDir, flag)
    await runExploit(i, workingDir, problem)
  }
}

const calcScore = async (team, problem) => {
  const flags = generateFlags(problem.exploit_trial_count)
  await runExploits(team, problem, flags)
  const submittedFlags = getSubmittedFlags(team, problem)

  const correctFlags = flags.filter(flag => submittedFlags.includes(flag))
  return correctFlags.length
}

module.exports = {
  calcScore
}
