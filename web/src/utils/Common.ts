import {AES, enc} from 'crypto-js';

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

export const isHTML = (text: string) => /^/.test(text);

export const getObjectIncludesText = (object: unknown, text: string = ''): boolean => {
  if (!text.length) return false;

  const searchTextLower = text.toLowerCase();
  const stringSpan = JSON.stringify(object).toLowerCase();

  return stringSpan.includes(searchTextLower);
};

export const downloadFile = (data: string, fileName: string): Element => {
  const element = document.createElement('a');
  const file = new Blob([data]);
  element.href = URL.createObjectURL(file);
  element.download = fileName;
  document.body.appendChild(element);
  element.click();
  return element;
};

export function enumKeys<O extends object, K extends keyof O = keyof O>(obj: O): K[] {
  return Object.keys(obj).filter(k => Number.isNaN(Number(k))) as K[];
}

export function singularOrPlural(noun: string, quantity: number) {
  if (quantity === 1) return noun;
  return `${noun}s`;
}

const encryptKey = 'tracetest';

export const encryptString = (data: string): string => {
  return AES.encrypt(data, encryptKey).toString();
};

export const decryptString = (data: string): string => {
  const bytes = AES.decrypt(data, encryptKey);

  return bytes.toString(enc.Utf8);
};
