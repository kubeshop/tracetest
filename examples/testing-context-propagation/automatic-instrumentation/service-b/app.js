const express = require("express")

const app = express()

const port = process.env.API_PORT || 8801

app.post("/augmentData", (req, res) => {
  const body = req.body

  const augmentedData = {
    ...body,
    "messageFromB": "Data augmented by service B"
  }

  return res.send(augmentedData)
})

app.listen(port, () => {
  console.log(`Listening for requests on http://localhost:${port}`)
})
