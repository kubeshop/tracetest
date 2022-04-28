export type TRecursivePartial<T> = {
  [P in keyof T]?: TRecursivePartial<T[P]>;
};

export type TEnv = {
  measurementId: string;
  analyticsEnabled: string;
};
