import LocalStorageGateway from 'gateways/LocalStorage.gateway';

const storageKey = 'user_preferences';

interface IUserPreferences {
  lang: string;
}

const localStorageGateway = LocalStorageGateway<IUserPreferences>(storageKey);

const initialUserPreferences: IUserPreferences = {
  lang: 'en',
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
