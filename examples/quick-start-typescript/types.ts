export type PokemonList = {
  items: Pokemon[];
  totalCount: number;
};

export type Pokemon = {
  id: number;
  imageUrl: string;
  isFeatured: boolean;
  type: string;
  name: string;
};
