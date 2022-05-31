import {render} from '../../../test-utils';
import {TSpanFlatAttribute} from '../../../types/Span.types';
import AttributeRow from '../AttributeRow';

const attribute: TSpanFlatAttribute = {
  key: 'key',
  value: 'value',
};

const onCreateAssertion = jest.fn();
const onCopy = jest.fn();
const setIsCopied = jest.fn();

describe('AttributeRow', () => {
  it('should render correctly', () => {
    const {getByText} = render(
      <AttributeRow
        attribute={attribute}
        isCopied={false}
        onCreateAssertion={onCreateAssertion}
        onCopy={onCopy}
        setIsCopied={setIsCopied}
      />
    );

    expect(getByText(attribute.key)).toBeInTheDocument();
    expect(getByText(attribute.value)).toBeInTheDocument();
  });
});
