const express = require("express")
const app = express()
app.get("/", (req, res) => {
  setTimeout(() => {
    res.send("Hello World")
  }, 1000);
})
app.listen(8084, () => {
  console.log(`Listening for requests on http://localhost:8080`)
})
