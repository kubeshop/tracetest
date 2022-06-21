// import {ReactFlowProvider} from 'react-flow-renderer';
import AssertionResultsMock from '../../../models/__mocks__/AssertionResults.mock';
import {render} from '../../../test-utils';
import AssertionCard from '../AssertionCard';

const onSelectSpan = jest.fn();
const onDelete = jest.fn();
const onEdit = jest.fn();

describe('AssertionCard', () => {
  it('should render', () => {
    const {
      resultList: [assertionResult],
    } = AssertionResultsMock.model();

    const {getByTestId} = render(
      // <ReactFlowProvider>
      <AssertionCard
        assertionResult={assertionResult}
        onDelete={onDelete}
        onEdit={onEdit}
        onSelectSpan={onSelectSpan}
      />
      // </ReactFlowProvider>
    );

    expect(getByTestId('assertion-card')).toBeInTheDocument();
  });
});
