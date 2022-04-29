export type TRecursivePartial<T> = {
  [P in keyof T]?: TRecursivePartial<T[P]>;
};

export interface IEnv {
  measurementId: string;
  analyticsEnabled: string;
}
