import {render} from '@testing-library/react';
import Request from '../components/Http/Request';

test('Request', () => {
  const {getAllByTestId} = render(
    <Request
      attributeList={[
        {key: 'http.user_agent', value: 'iPhone'},
        {key: 'http.request_content_length', value: '3024'},
        {key: 'http.not', value: 'nothing'},
        {key: 'http.host', value: 'google.com'},
      ]}
    />
  );
  expect(getAllByTestId('assertion-check-property').length).toBe(3);
});
