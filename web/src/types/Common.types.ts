export type TRecursivePartial<T> = {
  [P in keyof T]?: TRecursivePartial<T[P]>;
};

export interface IEnv {
  measurementId: string;
  analyticsEnabled: string;
};

export type Modify<T, R> = Omit<T, keyof R> & R;
