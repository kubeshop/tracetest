import {render} from '../../../test-utils';
import {TSpanFlatAttribute} from '../../../types/Span.types';
import AttributeRow from '../AttributeRow';

const attribute: TSpanFlatAttribute = {
  key: 'key',
  value: 'value',
};

const onCreateOutput = jest.fn();
const onCreateTestSpec = jest.fn();

describe('AttributeRow', () => {
  it('should render correctly', () => {
    const {getByText} = render(
      <AttributeRow
        searchText=""
        attribute={attribute}
        onCreateOutput={onCreateOutput}
        onCreateTestSpec={onCreateTestSpec}
        semanticConventions={{}}
        outputs={[]}
      />
    );

    expect(getByText(attribute.key)).toBeInTheDocument();
    expect(getByText(attribute.value)).toBeInTheDocument();
  });
});
