export interface IUserPreferences {
  lang: string;
  environmentId: string;
  initConfigSetup: boolean;
  initConfigSetupFromTest: boolean;
  isOnboardingComplete: boolean;
}

export type TUserPreferenceKey = keyof IUserPreferences;
export type TUserPreferenceValue<K extends TUserPreferenceKey = TUserPreferenceKey> = IUserPreferences[K];

export interface IUserState {
  preferences: IUserPreferences;
}
