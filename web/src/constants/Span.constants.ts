export const enum LOCATION_NAME {
  RESOURCE_ATTRIBUTES = 'RESOURCE_ATTRIBUTES',
  INSTRUMENTATION_LIBRARY = 'INSTRUMENTATION_LIBRARY',
  SPAN = 'SPAN',
  SPAN_ATTRIBUTES = 'SPAN_ATTRIBUTES',
  SPAN_ID = 'SPAN_ID',
}

export const HttpRequestAttributeList = [
  'http.method',
  'http.url',
  'http.target',
  'http.host',
  'http.scheme',
  'http.request_content_length',
  'http.request_content_length_uncompressed',
  'http.retry_count exists',
  'http.user_agent',
];

export const HttpResponseAttributeList = ['http.status_code'];
