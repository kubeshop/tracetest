/**
 * Welcome to Cloudflare Workers! This is your first worker.
 *
 * - Run `npm run dev` in your terminal to start a development server
 * - Open a browser tab at http://localhost:8787/ to see your worker in action
 * - Run `npm run deploy` to publish your worker
 *
 * Bind resources to your worker in `wrangler.toml`. After adding bindings, a type definition for the
 * `Env` object can be regenerated with `npm run cf-typegen`.
 *
 * Learn more at https://developers.cloudflare.com/workers/
 */

/**
 * STEP 1
 */

// export interface Env {
//   DB: D1Database
// }

// export async function addPokemon(pokemon: any, env: Env) {
//   return await env.DB.prepare(
//     "INSERT INTO Pokemon (name) VALUES (?) RETURNING *"
//   ).bind(pokemon.name).all()
// }

// async function formatPokeApiResponse(response: any) {
//   const data = await response.json()
//   const { name, id } = data
//   return { name, id }
// }

// export default {
// 	async fetch(request, env, ctx): Promise<Response> {
//     try {
//       const { pathname, searchParams } = new URL(request.url)

//       // Import a Pokemon
//       if (pathname === "/api/pokemon" && request.method === "POST") {
//         const response = await fetch(`https://pokeapi.co/api/v2/pokemon/${searchParams.get('id') || '6'}`)
//         const resPokemon = await formatPokeApiResponse(response)

//         const addedPokemon = await addPokemon(resPokemon, env)
//         return Response.json(addedPokemon)
//       }

//       return new Response("Hello Worker!")
//     } catch (err) {
//       return new Response(String(err))
//     }
// 	},
// } satisfies ExportedHandler<Env>;


/**
 * STEP 1 END
 */

/**
 * STEP 2
 */

import { trace, SpanStatusCode } from '@opentelemetry/api'
import { instrument, ResolveConfigFn } from '@microlabs/otel-cf-workers'
const tracer = trace.getTracer('pokemon-api')

export interface Env {
  DB: D1Database
  GRAFANA_URL: string
  GRAFANA_AUTH: string
}

// Manual OpenTelemetry instrumentation Start
export async function addPokemonWithTrace(pokemon: any, env: Env) {
  return tracer.startActiveSpan('D1 Add Pokemon', async (span) => {
    const addedPokemon = await env.DB.prepare(
      "INSERT INTO Pokemon (name) VALUES (?) RETURNING *"
    ).bind(pokemon.name).all()

    span.setStatus({
      code: SpanStatusCode.OK,
      message: String("Pokemon added successfully!")
    })
    span.setAttribute('pokemon.name', String(addedPokemon?.results[0].name))
    span.end()
    
    return addedPokemon
  })
}

async function formatPokeApiResponseWithTrace(response: any) {
  const span = trace.getActiveSpan()
  try {
    const { name, id } = await response.json()
    if (span) {
      span.setAttribute('pokeapi.response.name', String(name))
      span.setAttribute('pokeapi.response.id', String(id))
      span.setAttribute('pokeapi.response.error', Boolean(false))
    }
    return { name, id }
  } catch (err) {
    if (span) {
      span.setStatus({ code: SpanStatusCode.ERROR,
        message: String("Failed to fetch Pokemon.") })
      span.setAttribute('pokeapi.response.error', Boolean(true))
      span.setAttribute('pokeapi.response.error.message', String(err))
    }
    throw new Error(String(err))
  }
}
// Manual OpenTelemetry instrumentation End

const handler = {
	async fetch(request, env, ctx): Promise<Response> {
    try {
      const { pathname, searchParams } = new URL(request.url)

      // Import a Pokemon
      if (pathname === "/api/pokemon" && request.method === "POST") {
        const response = await fetch(`https://pokeapi.co/api/v2/pokemon/${searchParams.get('id') || '6'}`)
        const resPokemon = await formatPokeApiResponseWithTrace(response)
        const addedPokemon = await addPokemonWithTrace(resPokemon, env)
        return Response.json(addedPokemon)
      }

      return new Response("Hello Worker!")
    } catch (err) {
      return new Response(String(err))
    }
	},
} satisfies ExportedHandler<Env>;

const config: ResolveConfigFn = (env: Env, _trigger) => {
  return {
    exporter: {
      url: env.GRAFANA_URL,
      headers: { 'authorization': env.GRAFANA_AUTH },
    },
		service: { name: 'pokemon-api' },
	}
}

export default instrument(handler, config)

/**
 * STEP 2 END
 */
