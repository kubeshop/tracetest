import {capitalize} from 'lodash';
import Env from './Env';

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

export function singularOrPlural(noun: string, quantity: number) {
  if (quantity === 1) return noun;
  return `${noun}s`;
}

export function ordinalSuffixOf(i: number) {
  const j = i % 10;
  const k = i % 100;
  if (j === 1 && k !== 11) {
    return `${i}st`;
  }
  if (j === 2 && k !== 12) {
    return `${i}nd`;
  }
  if (j === 3 && k !== 13) {
    return `${i}rd`;
  }
  return `${i}th`;
}

export function getTotalCountFromHeaders(meta: any) {
  return Number(meta?.response?.headers.get('x-total-count') || 0);
}

export const getServerBaseUrl = () => {
  const {host, protocol} = window.location;
  const prefix = Env.get('serverPathPrefix');

  return `${protocol}//${host}${prefix}`;
};

export const ToTitle = (str: string) => {
  return capitalize(str.replace(/\W/g, ' '));
};

export const getIsValidUrl = (url: string): boolean => {
  try {
    return !!getParsedURL(url);
  } catch (e) {
    return false;
  }
};

export const getParsedURL = (rawUrl: string): URL => {
  if (!!rawUrl && !rawUrl.startsWith('http')) {
    return new URL(`http://${rawUrl}`);
  }

  return new URL(rawUrl);
};
