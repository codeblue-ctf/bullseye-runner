const uuidv4 = require('uuid/v4')
const config = require('../config.js')

const generateFlag = () => {
  const flagContent = uuidv4()
  return config.flagTemplate.replace('{flag}', flagContent)
}

const generateFlags = (count) => {
  const flags = []
  for (let i = 0; i < count; ++i) {
    flags.push(generateFlag())
  }
  return flags
}

module.exports = {
  generateFlags
}
