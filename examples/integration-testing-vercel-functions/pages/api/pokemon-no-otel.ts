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
  try {
    const requestUrl = `https://pokeapi.co/api/v2/pokemon/${req.body.id || '6'}`
    const response = await fetch(requestUrl)
    const resPokemon = await response.json()
    
    const { rowCount, rows: [addedPokemon, ...addedPokemonRest] } = await addPokemon(resPokemon) 
    res.status(200).json(addedPokemon)

  } catch (err) {
    res.status(500).json({ error: 'failed to load data' })
  }
}
