import {render} from '@testing-library/react';
import Response from '../components/Http/Response';

test('Response', () => {
  const {container} = render(
    <Response
      attributeList={[
        {key: 'http.response.token', value: '2309ei'},
        {key: 'http.status_code', value: '200'},
        {key: 'http.response.content', value: '21232'},
      ]}
    />
  );
  expect(container).toMatchSnapshot();
  expect(container.getElementsByClassName('ant-table-row').length).toBe(2);
});
