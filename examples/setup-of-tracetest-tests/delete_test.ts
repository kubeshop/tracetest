
import fs from 'fs';
import { Pokemon } from './types';
import Tracetest from '@tracetest/client';

// To use the @tracetest/client, you must have a token for the environment. This is created on the Settings 
// page under Tokens by the administrator for the environment. The token below has been given the 'engineer'
// role in the pokeshop-demo env in the tracetest-demo org so you can create and run tests in this environment.
// Want to read more about setting up tokens? https://docs.tracetest.io/concepts/environment-tokens

const TRACETEST_API_TOKEN = 'tttoken_4fea89f6e7fa1500';
const baseUrl = 'https://demo-pokeshop.tracetest.io/pokemon';

// json for the body of the POST to create a pokemon
const setupPokemon = `{
    "name": "fearow",  
    "imageUrl": "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/22.png",
    "isFeatured": false,
    "type": "normal,flying"
}`

const main = async () => {
    console.log('Lets use the TRACETEST_API_TOKEN to authenticate the @tracetest/client module...')
    const tracetest = await Tracetest(TRACETEST_API_TOKEN);
    

    //execute setup by adding a Pokemon with a REST POST api call directly
    const createPokemon = async (): Promise<Pokemon> => {
        const response = await fetch(baseUrl,{
            method: 'POST',
            body: setupPokemon,
            headers: {'Content-Type': 'application/json'} 
        });
        return await response.json() as Pokemon;
    };

    console.log('Adding the Pokemon - this is the setup action we need before running a Tracetest test')
    const pokemon = await createPokemon();

    // Get the id of the pokemon that we created in the setup step
    let pokemonId = pokemon.id;
    console.log('The Pokemon id we created was ', pokemonId);

    
    // Lets pull in the delete-pokemon test from a file
    let deleteTest = fs.readFileSync('delete_pokemon.yaml', 'utf-8');

    // Lets setup the variables we will be passing into the test (ie the pokemon_id)
    const getVariables = (id: string) => [
        { key: 'pokemon_id', value: id }
    ];

    const deletePokemon = async () => {
        console.log('Creating the delete-pokemon test based on the test in delete_pokemon.yaml...');
        const test = await tracetest.newTest(deleteTest);


        // run deletes pokemon test
        console.log('Running the delete-pokemon test...');
        const run = await tracetest.runTest(test, { variables: getVariables(String(pokemonId)) });
        await run.wait();
    };

    await deletePokemon();
    
    console.log(await tracetest.getSummary());
};

main();
