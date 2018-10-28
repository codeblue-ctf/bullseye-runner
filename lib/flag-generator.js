const uuidv4 = require('uuid/v4')

const generateFlag = (flagTemplate) => {
  const flagContent = uuidv4()
  return flagTemplate.replace('{{flag}}', flagContent)
}

const generateFlags = (count, flagTemplate) => {
  const flags = []
  for (let i = 0; i < count; ++i) {
    flags.push(generateFlag(flagTemplate))
  }
  return flags
}

module.exports = {
  generateFlags
}
