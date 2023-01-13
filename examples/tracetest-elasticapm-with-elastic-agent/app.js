// Add this to the VERY top of the first file loaded in your app
const apm = require('elastic-apm-node').start({
  serviceName: 'sample-app',
  serverUrl: 'http://apm-server:8200',
})

const express = require("express")
const app = express()
app.get("/", (req, res) => {
  res.send("Hello World")
})
app.listen(8080, () => {
  console.log(`Listening for requests on http://localhost:8080`)
})
