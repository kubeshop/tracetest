const { POKESHOP_DEMO_URL = '',  TYK_AUTH_KEY = '' } = process.env;

export const getKey = async () => {
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'x-tyk-authorization': TYK_AUTH_KEY,
      'Response-Type': 'application/json',
    },
  };

  const data = {
    alias: 'website',
    expires: -1,
    access_rights: {
      1: {
        api_id: '1',
        api_name: 'pokeshop',
        versions: ['Default'],
      },
    },
  };

  const res = await fetch(`${POKESHOP_DEMO_URL}/tyk/keys/create`, {
    ...params,
    method: 'POST',
    body: JSON.stringify(data),
  });

  const { key } = (await res.json()) as { key: string };

  return key;
};
