const { spawnSync } = require('child_proces')
const fs = require('fs')
const path = require('path')
const config = require('../config.js')
const { generateFlags } = require('./lib/flag-generator.js')

const loginRegistry = () => {
  spawnSync('docker', [
    'login',
    config.registry.server,
    '-u',
    config.registry.usre.name,
    '-p',
    config.registry.usre.password
  ])
}

const pullImage = (team, problem) => {
  const image = `${team.name}/${problem.name}`
  spawnSync('docker', ['pull', image])
  spawnSync('docker', ['tag', image, problem.exploit_container_name])
}

const setFlag = (workingDir, flag) => {
  // XXX: flag filename is hardcoded
  const flagFile = path.join(workingDir, 'flag')
  fs.writeFileSync(flagFile, flag)
}

const runExploit = (workingDir, problem) => {
  spawnSync('docker-compose', ['up'], { cwd: workingDir })
  // TODO: sleep problem.run_time
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
  const flags = generateFlags(problem.calc_time)
  runExploits(team, problem, flags)
  const submittedFlags = getSubmittedFlags(team, problem)

  const correctFlags = flags.filter(flag => submittedFlags.includes(flag))
  return correctFlags.length
}
