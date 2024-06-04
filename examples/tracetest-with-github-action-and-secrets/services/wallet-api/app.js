const express = require("express")
const database = require("./data.json")
const config = require("./config")

const app = express()

app.get("/wallet/:id", (req, res) => {
  const { id } = req.params
  const wallet = database.find(item => item.id == id)

  if (!wallet) {
    return res.sendStatus(404)
  }

  return res.send(wallet)
})

app.listen(config.port, () => {
  console.log(`Listening for requests on http://localhost:${config.port}`)
})
