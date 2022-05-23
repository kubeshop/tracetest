import SpanAttribute from '../SpanAttribute.model';
import SpanAttributeMock from '../__mocks__/SpanAttribute.mock';

describe('Span Attribute', () => {
  it('should generate a span attribute object', () => {
    const rawSpanAttribute = SpanAttributeMock.raw();
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual(rawSpanAttribute.value);
  });
});
