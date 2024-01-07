import { trace, SpanStatusCode } from '@opentelemetry/api'
import type { NextApiRequest, NextApiResponse } from 'next'
import { sql } from '@vercel/postgres'

export async function addPokemon(pokemon: any) {
  return await sql`
    INSERT INTO pokemon (name)
    VALUES (${pokemon.name})
    RETURNING *;
  `
}

export async function getPokemon(pokemon: any) {
  return await sql`
    SELECT * FROM pokemon where id=${pokemon.id};
  `
}

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const activeSpan = trace.getActiveSpan()
  const tracer = await trace.getTracer('integration-testing-vercel-functions')
  
  try {

    const externalPokemon = await tracer.startActiveSpan('GET Pokemon from pokeapi.co', async (externalPokemonSpan) => {
      const requestUrl = `https://pokeapi.co/api/v2/pokemon/${req.body.id || '6'}`
      const response = await fetch(requestUrl)
      const { id, name } = await response.json()

      externalPokemonSpan.setStatus({ code: SpanStatusCode.OK, message: String("Pokemon fetched successfully!") })
      externalPokemonSpan.setAttribute('pokemon.name', name)
      externalPokemonSpan.setAttribute('pokemon.id', id)
      externalPokemonSpan.end()

      return { id, name }
    })

    const addedPokemon = await tracer.startActiveSpan('Add Pokemon to Vercel Postgres', async (addedPokemonSpan) => {
      const { rowCount, rows: [addedPokemon, ...rest] } = await addPokemon(externalPokemon)
      addedPokemonSpan.setAttribute('pokemon.isAdded', rowCount === 1)
      addedPokemonSpan.setAttribute('pokemon.added.name', addedPokemon.name)
      addedPokemonSpan.end()
      return addedPokemon
    })
    
    res.status(200).json(addedPokemon)

  } catch (err) {
    activeSpan?.setAttribute('error', String(err))
    activeSpan?.recordException(String(err))
    activeSpan?.setStatus({ code: SpanStatusCode.ERROR, message: String(err) })
    res.status(500).json({ error: 'failed to load data' })
  } finally {
    activeSpan?.end()
  }
}
