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

  // TODO: build flag submit container
  spawnSyncWithLog('docker', ['build', 'flag-submit', '-t', config.flagSubmitContainerName])
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

  // clear submitted-flag
  // XXX: submitted-flag filename is hardcoded
  const submittedFlag = path.join(workingDir, 'submitted-flag')
  fs.writeFileSync(submittedFlag, '')
}

const runExploit = async (i, workingDir, problem) => {
  console.debug('start exploit', i)

  // XXX: It should be `spawn` instead of `spawnSync`?
  // Malicious exploit container may take times to start up
  spawnSyncWithLog('docker-compose', ['up', '-d'], { cwd: workingDir })

  await sleep(problem.exploit_timeout * 1000)

  spawnSyncWithLog('docker-compose', ['down', '-t', 0], { cwd: workingDir })

  console.debug('end exploit', i)
}

const getSubmittedFlags = (workingDir) => {
  const submittedFlag = path.join(workingDir, 'submitted-flag')
  // TODO: should trim the flag
  return fs.readFileSync(submittedFlag, '').toString()
}

const runExploits = async (team, problem, flags) => {
  const workingDir = fs.mkdtempSync(path.join('tmp', 'bulls-eye-runner-'))
  console.debug('workingDir', workingDir)

  loginRegistry()
  pullImage(team, problem)
  setupDockerCompose(workingDir, problem)

  const submittedFlags = []
  for (let i = 0; i < problem.exploit_trial_count; ++i) {
    const flag = flags[i]
    setFlag(workingDir, flag)

    await runExploit(i, workingDir, problem)

    const submittedFlag = getSubmittedFlags(workingDir)
    submittedFlags.push(submittedFlag)
  }
  return submittedFlags
}

const calcScore = async (team, problem) => {
  const flags = generateFlags(problem.exploit_trial_count)
  const submittedFlags = await runExploits(team, problem, flags)

  const correctFlags = flags.filter(flag => submittedFlags.includes(flag))
  return correctFlags.length
}

module.exports = {
  calcScore
}
