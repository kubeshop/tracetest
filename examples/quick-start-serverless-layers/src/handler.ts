import { APIGatewayEvent, Handler } from 'aws-lambda';
import fetch from 'node-fetch';
import { Pokemon, RawPokemon } from './types';
import DynamoDb from './dynamodb';

const Pokemon = (raw: RawPokemon): Pokemon => ({
  id: raw.id,
  name: raw.name,
  types: raw.types.map((type) => type.type.name),
  imageUrl: raw.sprites.front_default,
});

const getPokemon = async (id: string): Promise<Pokemon> => {
  const url = `https://pokeapi.co/api/v2/pokemon/${id}`;
  const res = await fetch(url);

  const raw = await res.json();

  return Pokemon(raw);
};

const insertPokemon = async (pokemon: Pokemon) => {
  await DynamoDb.put(pokemon);

  return DynamoDb.get<Pokemon>(pokemon.id);
};

type TBody = { id: string };

export const importPokemon: Handler<APIGatewayEvent> = async (event) => {
  console.log(event);

  const { id = '' } = JSON.parse(event.body || '') as TBody;

  try {
    const pokemon = await getPokemon(id);
    const result = await insertPokemon(pokemon);

    return {
      statusCode: 200,
      body: JSON.stringify(result),
    };
  } catch (error) {
    return {
      statusCode: 400,
      body: error.message,
    };
  }
};
