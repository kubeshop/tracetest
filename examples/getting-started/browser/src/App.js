import "./App.css";
import { useState, useEffect } from "react";

function App() {
  const [pokemon, setPokemon] = useState([]);
  useEffect(() => {
    fetch(`https://pokeapi.co/api/v2/pokemon?limit=5&offset=0`)
      .then((res) => {
        return res.json();
      })
      .then((data) => {
        console.log(data.results);
        setPokemon(data.results);
      })
      .catch((err) => {
        console.log(err);
      });
  }, []);
  return (
    <div>
      <h1>Pokemon</h1>
      {pokemon.map((pokemon) => (
        <li key={pokemon.url}>{pokemon.name}</li>
      ))}
    </div>
  );
}

export default App;
