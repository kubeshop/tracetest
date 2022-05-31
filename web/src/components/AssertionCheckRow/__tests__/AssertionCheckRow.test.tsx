import AssertionMock from '../../../models/__mocks__/Assertion.mock';
import AssertionSpanResultMock from '../../../models/__mocks__/AssertionSpanResult.mock';
import {render} from '../../../test-utils';
import AssertionCheckRow from '../AssertionCheckRow';

const onSelectSpan = jest.fn();
const getIsSelectedSpan = jest.fn();

describe('AssertionCheckRow', () => {
  it('should render', () => {
    const result = AssertionSpanResultMock.model();
    const assertion = AssertionMock.model();

    const {getByText} = render(
      <AssertionCheckRow
        result={result}
        assertion={assertion}
        getIsSelectedSpan={getIsSelectedSpan}
        onSelectSpan={onSelectSpan}
      />
    );

    expect(getByText(assertion.attribute)).toBeInTheDocument();
  });
});
