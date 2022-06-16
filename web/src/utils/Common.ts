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

export const downloadFile = (data: string, fileName: string): void => {
  const element = document.createElement('a');
  const file = new Blob([data]);
  element.href = URL.createObjectURL(file);
  element.download = fileName;
  document.body.appendChild(element);
  element.click();
};
