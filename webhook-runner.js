const { calcScore } = require('./lib/calc-score')

process.on('message', async (data) => {
  const { id, team, problem } = data
  const score = await calcScore(team, problem)
  // TODO: callback
})
