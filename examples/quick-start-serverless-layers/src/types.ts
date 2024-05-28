export type RawPokemon = {
  id: number;
  name: string;
  types: Array<{
    type: {
      name: string;
    };
  }>;
  sprites: {
    front_default: string;
  };
};


export type Pokemon = {
  id: number;
  name: string;
  types: string[];
  imageUrl: string;
};
