module.exports = {
  flagTemplate: process.env.BULLSEYE_FLAG_TEMPLATE || 'CBCTF{{flag}}',
  flagSubmitContainerName: process.env.BULLSEYE_FLAG_SUBMIT_CONTAINER_NAME || 'flag-submit',
  registry: {
    server: process.env.BULLSEYE_REGISTRY_SERVER || 'localhost:5000',
    admin: {
      name: proces.env.BULLSEYE_REGISTRY_ADMIN_NAME || 'admin',
      password: process.env.BULLSEYE_REGISTRY_ADMIN_PASSWORD || 'password'
    }
  },
  bullseyeWeb: {
    server: process.env.BULLSEYE_WEB_SERVER || 'web:3000'
  }
}
