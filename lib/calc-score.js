const { spawnSync } = require('child_process')
const fs = require('fs')
const path = require('path')
const { generateFlags } = require('./flag-generator.js')
const { sleep } = require('./sleep.js')

const spawnSyncWithLog = (command, args, options) => {
  console.debug('[spawnSync]', command, args)
  const result = spawnSync(command, args, options)
  console.debug('[spawnSync][stdout]', result.stdout.toString())
  console.debug('[spawnSync][stderr]', result.stderr.toString())
  return result
}

const loginRegistry = (registryUrl, username, password) => {
  spawnSyncWithLog('docker', [
    'login',
    registryUrl,
    '-u',
    username,
    '-p',
    password
  ])
}

const setupDockerCompose = (workingDir, dockerCompose) => {
  const dockerComposeFile = path.join(workingDir, 'docker-compose.yml')
  fs.writeFileSync(dockerComposeFile, dockerCompose)

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

const runExploit = async (i, workingDir, timeout) => {
  console.debug('start exploit', i)

  // XXX: It should be `spawn` instead of `spawnSync`?
  // Malicious exploit container may take times to start up
  spawnSyncWithLog('docker-compose', ['up', '-d'], { cwd: workingDir })

  await sleep(timeout)

  spawnSyncWithLog('docker-compose', ['down', '-t', 0], { cwd: workingDir })

  console.debug('end exploit', i)
}

const getSubmittedFlags = (workingDir) => {
  const submittedFlag = path.join(workingDir, 'submitted-flag')
  // TODO: should trim the flag
  return fs.readFileSync(submittedFlag, '').toString()
}

const runExploits = async (config, flags) => {
  const workingDir = fs.mkdtempSync(path.join('tmp', 'bullseye-runner-'))
  console.debug('workingDir', workingDir)

  // login to docker registry and pull images
  loginRegistry(config.registry_host, config.admin_username, config.admin_password)

  setupDockerCompose(workingDir, config.docker_compose)

  const submittedFlags = []
  for (let i = 0; i < config.trials_count; ++i) {
    const flag = flags[i]
    setFlag(workingDir, flag)

    await runExploit(i, workingDir, config.timeout)

    const submittedFlag = getSubmittedFlags(workingDir)
    submittedFlags.push(submittedFlag)
  }
  return submittedFlags
}

/*
config = {
  trials_count,
  timeout,
  docker_compose,
  registry_host,
  admin_username,
  admin_password,
  flag_template
}
*/
const calcScore = async (config) => {
  const flags = generateFlags(config.trials_count, config.flag_template)
  const submittedFlags = await runExploits(config, flags)

  // TODO: return arrays which describe the exploit succeeded in that time
  const correctFlags = flags.filter(flag => submittedFlags.includes(flag))
  return correctFlags.length
}

module.exports = {
  calcScore
}
