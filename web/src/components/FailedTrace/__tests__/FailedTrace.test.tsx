import {render} from 'test-utils';
import FailedTrace from '../index';
import TestMock from '../../../models/__mocks__/Test.mock';
import TestRunMock from '../../../models/__mocks__/TestRun.mock';

test('FailedTrace', () => {
  const test = TestMock.model();
  const {getByText} = render(<FailedTrace testId={test.id} isDisplayingError run={TestRunMock.model()} />);
  expect(getByText('Rerun Test')).toBeTruthy();
});
