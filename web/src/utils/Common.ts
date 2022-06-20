export const escapeString = (str: string): string => {
  // eslint-disable-next-line no-control-regex
  return str.replace(/[\\"']/g, '\\$&').replace(/\u0000/g, '\\0');
};

export const isBoolean = (value: string): boolean => value === 'true' || value === 'false';

export const isJson = (str: string) => {
  try {
    JSON.parse(str);
  } catch (e) {
    return false;
  }

  return Number.isNaN(Number(str)) && !isBoolean(str) && true;
};

export const getObjectIncludesText = (object: unknown, text: string): boolean => {
  if (!text.length) return false;

  const searchTextLower = text.toLowerCase();
  const stringSpan = JSON.stringify(object).toLowerCase();

  return stringSpan.includes(searchTextLower);
};
