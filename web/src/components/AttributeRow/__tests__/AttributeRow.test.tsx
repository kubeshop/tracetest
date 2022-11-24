import {render} from '../../../test-utils';
import {TSpanFlatAttribute} from '../../../types/Span.types';
import AttributeRow from '../AttributeRow';

const attribute: TSpanFlatAttribute = {
  key: 'key',
  value: 'value',
};

const onCreateOutput = jest.fn();
const onCreateTestSpec = jest.fn();
const onCopy = jest.fn();

describe('AttributeRow', () => {
  it('should render correctly', () => {
    const {getByText} = render(
      <AttributeRow
        searchText=""
        attribute={attribute}
        onCreateOutput={onCreateOutput}
        onCreateTestSpec={onCreateTestSpec}
        onCopy={onCopy}
        semanticConventions={{}}
      />
    );

    expect(getByText(attribute.key)).toBeInTheDocument();
    expect(getByText(attribute.value)).toBeInTheDocument();
  });
});
