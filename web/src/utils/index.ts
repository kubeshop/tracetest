const flattenAttributesSelector = () =>
  `resourceSpans[].[instrumentationLibrarySpans[].spans[].attributes[].{key:key,value:value.*|[0]},resource.attributes[].{key:key,value: value.*|[0]}]|[][]`;

export const filterBySpanId = (spanId: string = '') =>
  `resourceSpans[?instrumentationLibrarySpans[?spans[?starts_with(spanId,'${spanId}')]]] | [].[instrumentationLibrarySpans[].spans[].attributes[].{key:key,value:value.*|[0],type:'span'},resource.attributes[].{key:key,value: value.*|[0],type:'resource'}]|[][]`;
