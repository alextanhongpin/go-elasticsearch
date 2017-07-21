
const consul = require('consul')()
const elasticsearch = require('elasticsearch')

consul.catalog.service.nodes('global-elastichsearch-check', (err, results) => {
  if (err) {
    throw err
  }
  if (!results.length) {
    throw new Error('No services found')
  }
  // Pick a random service
  const result = results[Math.floor(Math.random() * results.length)]

  const client = new elasticsearch.Client({
    // host: 'http://elastic:changeme@127.0.0.1:20575',
    log: 'trace',
    hosts: [
      {
        host: result.ServiceAddress,
        auth: 'elastic:changeme',
        protocol: 'http',
        port: result.ServicePort
      }
    ]
  })

  client.ping({
    // ping usually has a 3000ms timeout
    requestTimeout: 1000
  }, (error) => {
    if (error) {
      throw error
    } else {
      console.log('Success')
    }
  })
})
