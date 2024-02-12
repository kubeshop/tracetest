import { trace, SpanStatusCode, diag, DiagConsoleLogger, DiagLogLevel } from '@opentelemetry/api'
import { instrument, ResolveConfigFn } from '@microlabs/otel-cf-workers'
const tracer = trace.getTracer('pokemon-api')

diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.DEBUG)

export interface Env {
  DB: D1Database
	TRACETEST_URL: string
}

export async function addPokemon(pokemon: any, env: Env) {
  return await env.DB.prepare(
    "INSERT INTO Pokemon (name) VALUES (?) RETURNING *"
  ).bind(pokemon.name).all()
}

export async function getPokemon(pokemon: any, env: Env) {
  return await env.DB.prepare(
    "SELECT * FROM Pokemon WHERE id = ?;"
  ).bind(pokemon.id).all();
}

async function formatPokeApiResponse(response: any) {
  const { headers } = response
  const contentType = headers.get("content-type") || ""
  if (contentType.includes("application/json")) {
    const data = await response.json()
    const { name, id } = data

    // Add manual instrumentation
    const span = trace.getActiveSpan()
    if(span) {
      span.setStatus({ code: SpanStatusCode.OK, message: String("Pokemon fetched successfully!") })
      span.setAttribute('pokemon.name', name)
      span.setAttribute('pokemon.id', id)
    }

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

        // Add manual instrumentation
        return tracer.startActiveSpan('D1: Add Pokemon', async (span) => {
          const addedPokemon = await addPokemon(resPokemon, env)

          span.setStatus({ code: SpanStatusCode.OK, message: String("Pokemon added successfully!") })
          span.setAttribute('pokemon.name', String(addedPokemon?.results[0].name))
          span.end()
          
          return Response.json(addedPokemon)
        })
      }

      return new Response("Hello Worker!")
    } catch (err) {
      return new Response(String(err))
    }
	},
}

const config: ResolveConfigFn = (env: Env, _trigger) => {
  return {
    exporter: {
      url: env.TRACETEST_URL,
      headers: { },
    },
		service: { name: 'pokemon-api' },
	}
}

export default instrument(handler, config)
