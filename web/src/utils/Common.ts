export const filterBySpanId = (spanId: string = '') =>
  `resourceSpans[?instrumentationLibrarySpans[?spans[?starts_with(spanId,'${spanId}')]]] | [].[instrumentationLibrarySpans[].spans[].attributes[].{key:key,value:value.*|[0],type:'span'},resource.attributes[].{key:key,value: value.*|[0],type:'resource'}]|[][]`;

export const escapeString = (str: string): string => {
  // eslint-disable-next-line no-control-regex
  return str.replace(/[\\"']/g, '\\$&').replace(/\u0000/g, '\\0');
};

export const isJson = (str: string) => {
  try {
    JSON.parse(str);
  } catch (e) {
    return false;
  }
  return true;
};

const visiblePortion = 100;

export function visiblePortionFuction() {
  return {visiblePortion, height: `calc(100% - ${visiblePortion}px - 77px)`};
}
