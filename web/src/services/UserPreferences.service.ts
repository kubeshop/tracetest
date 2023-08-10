import LocalStorageGateway from 'gateways/LocalStorage.gateway';
import {IUserPreferences} from 'types/User.types';

const storageKey = 'user_preferences';

const localStorageGateway = LocalStorageGateway<IUserPreferences>(storageKey);

const initialUserPreferences: IUserPreferences = {
  lang: 'en',
  variableSetId: '',
  initConfigSetup: true,
  initConfigSetupFromTest: true,
  showGuidedTourNotification: true,
  showAttributeTooltip: true,
};

const UserPreferencesService = () => ({
  get(): IUserPreferences {
    const userPreferences = localStorageGateway.get() || initialUserPreferences;

    return {
      ...initialUserPreferences,
      ...userPreferences,
    };
  },
  getEntry<K extends keyof IUserPreferences>(key: K): IUserPreferences[K] {
    const preferences = this.get();

    return preferences[key];
  },
  set<K extends keyof IUserPreferences>(key: K, value: IUserPreferences[K]) {
    const preferences = this.get();

    const updatedUserPreferences = {
      ...preferences,
      [key]: value,
    };

    localStorageGateway.set(updatedUserPreferences);

    return updatedUserPreferences;
  },
});

export default UserPreferencesService();
