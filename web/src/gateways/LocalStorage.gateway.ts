const LocalStorageGateway = <T>(defaultKey = '') => {
  const localstorage = window.localStorage;

  return {
    set<K = T>(value: K, key = defaultKey): void {
      const json = JSON.stringify(value);

      localstorage.setItem(key, json);
    },
    get<K = T>(key = defaultKey): K | undefined {
      const value = localstorage.getItem(key);

      return value ? JSON.parse(value) : undefined;
    },
  };
};

export default LocalStorageGateway;
