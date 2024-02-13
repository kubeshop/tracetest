export interface Env {
  DB: D1Database
}

export async function addPokemon(pokemon: any, env: Env) {
  return await env.DB.prepare(
    "INSERT INTO Pokemon (name) VALUES (?) RETURNING *"
  ).bind(pokemon.name).all()
}

export async function getPokemon(pokemon: any, env: Env) {
  return await env.DB.prepare(
    "SELECT * FROM Pokemon WHERE id = ?"
  ).bind(pokemon.id).all()
}

async function formatPokeApiResponse(response: any) {
  const { headers } = response
  const contentType = headers.get("content-type") || ""
  if (contentType.includes("application/json")) {
    const data = await response.json()
    const { name, id } = data
    return { name, id }
  }
  return response.text()
}

const handler = {
	async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    try {
      const { pathname, searchParams } = new URL(request.url)

      // Import a Pokemon
      if (pathname === "/api/pokemon" && request.method === "POST") {
        const queryId = searchParams.get('id')
        const requestUrl = `https://pokeapi.co/api/v2/pokemon/${queryId || '6'}`
        const response = await fetch(requestUrl)
        const resPokemon = await formatPokeApiResponse(response)

        const addedPokemon = await addPokemon(resPokemon, env)
        return Response.json(addedPokemon)
      }

      return new Response("Hello Worker!")
    } catch (err) {
      return new Response(String(err))
    }
	},
}

export default handler
