import {render} from 'test-utils';
import FailedTrace from '../index';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';

test('FailedTrace', () => {
  const {getByText} = render(<FailedTrace isDisplayingError run={TestRunMock.model()} />);
  expect(getByText('Test Run Failed')).toBeTruthy();
});
