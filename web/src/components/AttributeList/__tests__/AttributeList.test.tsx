import {render} from '../../../test-utils';
import AttributeList from '../AttributeList';

const onCreateOutput = jest.fn();
const onCreateTestSpec = jest.fn();

describe('AttributeList', () => {
  it('should render correctly', () => {
    const attributeList = [
      {
        key: 'key',
        value: 'value',
      },
    ];

    const {getByTestId} = render(
      <AttributeList
        attributeList={attributeList}
        onCreateOutput={onCreateOutput}
        onCreateTestSpec={onCreateTestSpec}
        semanticConventions={{}}
        outputs={[]}
      />
    );

    expect(getByTestId('attribute-list')).toBeInTheDocument();
  });

  it('should render the empty list', () => {
    const {getByTestId} = render(
      <AttributeList
        attributeList={[]}
        onCreateOutput={onCreateOutput}
        onCreateTestSpec={onCreateTestSpec}
        semanticConventions={{}}
        outputs={[]}
      />
    );

    expect(getByTestId('empty-attribute-list')).toBeInTheDocument();
  });
});
