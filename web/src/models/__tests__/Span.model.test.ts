import Span from '../Span.model';
import SpanMock from '../__mocks__/Span.mock';

describe('Span', () => {
  it('should generate a span object', () => {
    const rawSpan = SpanMock.raw();
    const span = Span(rawSpan);

    expect(Array.isArray(span.signature)).toBe(true);
    expect(typeof span.attributes).toEqual('object');
    expect(Array.isArray(span.attributeList)).toBe(true);
    const attribute = rawSpan.attributes[0];
    expect(span.attributes[attribute.key].value).toEqual(attribute.value.stringValue);

    const duration = Number(
      ((Number(rawSpan.endTimeUnixNano) - Number(rawSpan.startTimeUnixNano)) / 1000 / 1000).toFixed(1)
    );

    expect(span.duration).toEqual(duration);
  });

  describe('createFromResourceSpanList', () => {
    it('should return a list of spans from a resource span list', () => {
      const firstRawSpan = SpanMock.raw();
      const secondRawSpan = SpanMock.raw();
      const instrumentationLibrary = {
        name: 'test',
        version: '1',
      };
      const resource = {
        attributes: [],
      };

      const spanList = Span.createFromResourceSpanList([
        {
          resource,
          instrumentationLibrarySpans: [
            {
              instrumentationLibrary,
              spans: [firstRawSpan],
            },
          ],
        },
        {
          resource,
          instrumentationLibrarySpans: [
            {
              instrumentationLibrary,
              spans: [secondRawSpan],
            },
          ],
        },
      ]);

      expect(spanList).toHaveLength(2);
      const [spanOne, spanTwo] = spanList;

      expect(spanOne.spanId).toEqual(firstRawSpan.spanId);
      expect(spanTwo.spanId).toEqual(secondRawSpan.spanId);
    });

    it('should handle an empty list of resource spans', () => {
      const spanList = Span.createFromResourceSpanList([]);
  
      expect(spanList).toHaveLength(0);
    });
  });
});
