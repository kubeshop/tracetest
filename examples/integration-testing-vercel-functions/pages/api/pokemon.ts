import { trace, SpanStatusCode, Exception } from '@opentelemetry/api'
import type { NextApiRequest, NextApiResponse } from 'next'
 
export async function getTracer() {
  return await trace.getTracer('next-app')
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const tracer = await getTracer()
  tracer.startActiveSpan('GET Pokemon API', async (span) => {
    try {
      const requestUrl = `https://pokeapi.co/api/v2/pokemon/${req.query.id || '6'}`
      const response = await fetch(requestUrl)
      const data = await response.json()

      span.setStatus({ code: SpanStatusCode.OK, message: String("Pokemon fetched successfully!") })
      span.setAttribute('pokemon.name', data.name)
      span.setAttribute('pokemon.id', data.id)

      res.status(200).json({
        name: data.name,
      })
  
    } catch (err) {
      span.setAttribute('error', String(err))
      span.recordException(String(err))
      span.setStatus({ code: SpanStatusCode.ERROR, message: String(err) })
      res.status(500).json({ error: 'failed to load data' })
    } finally {
      span.end()
    }
  })
}
