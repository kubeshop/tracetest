import LocalStorageGateway from '../gateways/LocalStorage.gateway';
import {TViewResultsMode} from '../types/TestDefinition.types';

const storageKey = 'user_preferences';

interface IUserPreferences {
  viewResultsMode: TViewResultsMode;
}

const localStorageGateway = LocalStorageGateway<IUserPreferences>(storageKey);

const initialUserPreferences: IUserPreferences = {
  viewResultsMode: 'wizard',
};

const UserPreferencesService = () => ({
  getUserPreferences(): IUserPreferences {
    const userPreferences = localStorageGateway.get() || initialUserPreferences;

    return userPreferences;
  },
  getUserPreference<K extends keyof IUserPreferences>(key: K): IUserPreferences[K] {
    const preferences = this.getUserPreferences();

    return preferences[key];
  },
  setPreference<K extends keyof IUserPreferences>(key: K, value: IUserPreferences[K]) {
    const preferences = this.getUserPreferences();

    const updatedUserPreferences = {
      ...preferences,
      [key]: value,
    };

    localStorageGateway.set(updatedUserPreferences);

    return updatedUserPreferences;
  },
});

export default UserPreferencesService();
