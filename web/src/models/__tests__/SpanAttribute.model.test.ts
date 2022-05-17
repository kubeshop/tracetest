import SpanAttribute from '../SpanAttribute.model';
import SpanAttributeMock from '../__mocks__/SpanAttribute.mock';

describe('Span Attribute', () => {
  it('should generate a span attribute object', () => {
    const rawSpanAttribute = SpanAttributeMock.raw();
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual(rawSpanAttribute.value.stringValue);
    expect(spanAttribute.name).toEqual(rawSpanAttribute.key);
    expect(spanAttribute.type).toEqual('stringValue');
  });

  it('should generate a span attribute object with double value', () => {
    const rawSpanAttribute = SpanAttributeMock.raw({
      value: {
        doubleValue: 1.1,
        kvlistValue: {values: []},
      },
    });
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual(String(rawSpanAttribute.value.doubleValue));
    expect(spanAttribute.type).toEqual('doubleValue');
  });

  it('should generate a span attribute object with boolean value', () => {
    const rawSpanAttribute = SpanAttributeMock.raw({
      value: {
        booleanValue: true,
        kvlistValue: {values: []},
      },
    });
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual(String(rawSpanAttribute.value.booleanValue));
    expect(spanAttribute.type).toEqual('booleanValue');
  });

  it('should generate a span attribute object with a truthy boolean value', () => {
    const rawSpanAttribute = SpanAttributeMock.raw({
      value: {
        booleanValue: true,
        kvlistValue: {values: []},
      },
    });
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual(String(rawSpanAttribute.value.booleanValue));
    expect(spanAttribute.type).toEqual('booleanValue');
  });

  it('should generate a span attribute object with a falsy boolean value', () => {
    const rawSpanAttribute = SpanAttributeMock.raw({
      value: {
        booleanValue: false,
        kvlistValue: {values: []},
      },
    });
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual(String(rawSpanAttribute.value.booleanValue));
    expect(spanAttribute.type).toEqual('booleanValue');
  });

  it('should generate a span attribute object with a int value', () => {
    const rawSpanAttribute = SpanAttributeMock.raw({
      value: {
        intValue: 10,
        kvlistValue: {values: []},
      },
    });
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual(String(rawSpanAttribute.value.intValue));
    expect(spanAttribute.type).toEqual('intValue');
  });

  it('should generate a span attribute object with a zero int value', () => {
    const rawSpanAttribute = SpanAttributeMock.raw({
      value: {
        intValue: 0,
        kvlistValue: {values: []},
      },
    });
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual(String(rawSpanAttribute.value.intValue));
    expect(spanAttribute.type).toEqual('intValue');
  });

  it('should generate a span attribute object with an empty string', () => {
    const rawSpanAttribute = SpanAttributeMock.raw({
      value: {
        stringValue: '',
        kvlistValue: {values: []},
      },
    });
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual('<Empty value>');
    expect(spanAttribute.type).toEqual('stringValue');
  });

  it('should generate a span attribute object with a kvlistValue', () => {
    const rawSpanAttribute = SpanAttributeMock.raw({
      value: {
        kvlistValue: {values: []},
      },
    });
    const spanAttribute = SpanAttribute(rawSpanAttribute);

    expect(spanAttribute.value).toEqual(JSON.stringify(rawSpanAttribute.value.kvlistValue));
    expect(spanAttribute.type).toEqual('kvlistValue');
  });
});
