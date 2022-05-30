import {render} from '../../../test-utils';
import AttributeList from '../AttributeList';

const onCreateAssertion = jest.fn();

describe('AttributeList', () => {
  it('should render correctly', () => {
    const attributeList = [
      {
        key: 'key',
        value: 'value',
      },
    ];

    const {getByTestId} = render(<AttributeList attributeList={attributeList} onCreateAssertion={onCreateAssertion} />);

    expect(getByTestId('attribute-list')).toBeInTheDocument();
  });

  it('should render the empty list', () => {
    const {getByTestId} = render(<AttributeList attributeList={[]} onCreateAssertion={onCreateAssertion} />);

    expect(getByTestId('empty-attribute-list')).toBeInTheDocument();
  });
});
