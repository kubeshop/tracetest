const express = require("express")
const axios = require("axios")
const config = require("./config")

const opentelemetry = require('@opentelemetry/api')

const app = express()
app.use(express.json())

app.post("/executePaymentOrder", async (req, res) => {
  const { walletId, yearsAsACustomer } = req.body

  try {
    const balance = await getWalletBalance(walletId)
    const status = await executePayment(balance, yearsAsACustomer)

    return res.send({ status })
  } catch (error) {
    const activeSpan = opentelemetry.trace.getActiveSpan()

    activeSpan.recordException(error)
    activeSpan.setStatus({ code: opentelemetry.SpanStatusCode.ERROR })

    console.error(error.stack)
    return res.status(500).send('internal error!')
  }
})

app.listen(config.port, () => {
  console.log(`Listening for requests on http://localhost:${config.port}`)
})

async function getWalletBalance(walletId) {
  const axiosConfig = {
    method: 'get',
    url: config.walletAPIEndpoint + "/" + walletId,
    headers: {
      'Content-Type': 'application/json'
    }
  }

  const response = await axios.request(axiosConfig)
  return response.data.balance
}

async function executePayment(balance, age) {
  const raw = JSON.stringify({
    "amount": balance,
    "age": age
  })

  const axiosConfig = {
    method: 'post',
    url: config.paymentExecutorAPIEndpoint,
    headers: {
      'Content-Type': 'application/json'
    },
    data: raw
  }

  const response = await axios.request(axiosConfig)
  return response.data.status
}
