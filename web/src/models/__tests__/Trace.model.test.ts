import Trace from '../Trace.model';
import TraceMock from '../__mocks__/Trace.mock';

describe('Trace', () => {
  it('should generate a trace object', () => {
    const rawTrace = TraceMock.raw();
    const trace = Trace(rawTrace);

    expect(trace.description).toEqual(rawTrace.description);
    const length = rawTrace.resourceSpans![0].instrumentationLibrarySpans[0].spans.length;

    expect(trace.spans).toHaveLength(length);
  });

  it('should handle empty values', () => {
    const rawTrace = TraceMock.raw({
      description: undefined,
      resourceSpans: undefined,
    });
    const trace = Trace(rawTrace);

    expect(trace.description).toEqual('');
    expect(trace.spans).toHaveLength(0);
  });
});
