// type definition needed for the json returned from creating a new pokemon

export type Pokemon = {
    id: number;
    imageUrl: string;
    isFeatured: boolean;
    type: string;
    name: string;
  };