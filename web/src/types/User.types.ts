export interface IUserPreferences {
  lang: string;
  environmentId: string;
  initConfigSetup: boolean;
}

export type TUserPreferenceKey = keyof IUserPreferences;
export type TUserPreferenceValue<K extends TUserPreferenceKey = TUserPreferenceKey> = IUserPreferences[K];

export interface IUserState {
  preferences: IUserPreferences;
}
