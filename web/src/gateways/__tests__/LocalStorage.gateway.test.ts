import LocalStorageGateway from '../LocalStorage.gateway';

const localStorage = LocalStorageGateway<{name: string}>('key');

describe('LocalStorageGateway', () => {
  it('should set a value and retrieve it afterwards', async () => {
    expect.assertions(1);
    const item = {name: 'test'};
    localStorage.set(item);

    const value = localStorage.get();

    expect(value).toEqual(item);
  });

  it('should return undefined as key is not found', async () => {
    const value = localStorage.get('not-found');

    expect(value).toEqual(undefined);
  });
});
