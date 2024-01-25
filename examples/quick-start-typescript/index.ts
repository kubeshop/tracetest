import Tracetest from '@tracetest/client';
import { config } from 'dotenv';
import { PokemonList } from './types';
import { deleteDefinition, importDefinition } from './definitions';

config();

const { TRACETEST_API_TOKEN = '', POKESHOP_DEMO_URL = 'http://api:8081' } = process.env;

const baseUrl = `${POKESHOP_DEMO_URL}/pokemon`;

const main = async () => {
  const tracetest = await Tracetest(TRACETEST_API_TOKEN, 'https://app-stage.tracetest.io', '');

  const getLastPokemonId = async (): Promise<number> => {
    const response = await fetch('http://localhost:8081/pokemon');
    const list = (await response.json()) as PokemonList;

    return list.items.length + 1;
  };

  // get the initial pokemon from the API
  const pokemonId = (await getLastPokemonId()) + 1;

  const getVariables = (id: string) => [
    { key: 'POKEMON_ID', value: id },
    { key: 'BASE_URL', value: baseUrl },
  ];

  const importedPokemonList: string[] = [];

  const importPokemons = async (startId: number) => {
    const test = await tracetest.newTest(importDefinition);
    // imports all pokemons
    await Promise.all(
      new Array(5).fill(0).map(async (_, index) => {
        console.log(`ℹ Importing pokemon ${startId + index + 1}`);
        const run = await tracetest.runTest(test, { variables: getVariables(`${startId + index + 1}`) });
        const updatedRun = await run.wait();
        const pokemonId = updatedRun.outputs?.find((output) => output.name === 'DATABASE_POKEMON_ID')?.value || '';

        console.log(`ℹ Adding pokemon ${pokemonId} to the list`);
        importedPokemonList.push(pokemonId);
      })
    );
  };

  const deletePokemons = async () => {
    const test = await tracetest.newTest(deleteDefinition);
    // deletes all pokemons
    await Promise.all(
      importedPokemonList.map(async (pokemonId) => {
        console.log(`ℹ Deleting pokemon ${pokemonId}`);
        const run = await tracetest.runTest(test, { variables: getVariables(pokemonId) });
        run.wait();
      })
    );
  };

  await importPokemons(pokemonId);
  console.log(await tracetest.getSummary());

  await deletePokemons();
  console.log(await tracetest.getSummary());
};

main();
