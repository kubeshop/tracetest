export const ERROR_HEADER: Record<string, string> = {
  attribute_naming_error: 'The following attributes do not adhere to the naming convention:',
  empty_attribute_error: 'The following attributes are empty:',
  missing_attribute_error: 'This span is missing the following required attributes:',
  ip_address_error: 'The following attributes are using IP addresses instead of DNS:',
  insecure_protocol_error: 'The following attributes are using insecure protocols:',
  api_key_leak_error: 'The following attributes are exposing API keys:',
};
