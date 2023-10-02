export interface IUserPreferences {
  lang: string;
  variableSetId: string;
  initConfigSetup: boolean;
  initConfigSetupFromTest: boolean;
  showGuidedTourNotification: boolean;
  showAttributeTooltip: boolean;
}

export type TUserPreferenceKey = keyof IUserPreferences;
export type TUserPreferenceValue<K extends TUserPreferenceKey = TUserPreferenceKey> = IUserPreferences[K];

export interface IUserState {
  preferences: IUserPreferences;
  runOriginPath: string;
}
