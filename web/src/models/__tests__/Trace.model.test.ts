import Trace from '../Trace.model';
import TraceMock from '../__mocks__/Trace.mock';

describe('Trace', () => {
  it('should generate a trace object', () => {
    const rawTrace = TraceMock.raw();
    const trace = Trace(rawTrace);

    expect(trace.traceId).toEqual(rawTrace.traceId);
    const length = Object.keys(rawTrace.flat!).length;

    expect(trace.spans).toHaveLength(length);
  });

  it('should handle empty values', () => {
    const rawTrace = TraceMock.raw({
      traceId: undefined,
      flat: undefined,
    });
    const trace = Trace(rawTrace);

    expect(trace.traceId).toEqual('');
    expect(trace.spans).toHaveLength(0);
  });
});
