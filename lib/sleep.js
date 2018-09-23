const sleep = (millsec) => {
  return new Promise((resolve, reject) => {
    setTimeout(resolve, millsec)
  })
}

module.exports = {
  sleep
}
