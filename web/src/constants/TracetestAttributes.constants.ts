/**
 * Moved from JSON files to TS files, as JSONs are not working in composite projects.
 * @see {@link https://github.com/TypeStrong/ts-loader/issues/905}
 */
export default {
  name: {
    description:
      'The span name concisely identifies the work represented by the Span, for example, an RPC method name, a function name, or the name of a subtask or stage within a larger computation. The span name SHOULD be the most general string that identifies a (statistically) interesting class of Spans, rather than individual Span instances while still being human-readable.',
    note: '',
    tags: ['id'],
  },
  'tracetest.span.duration': {
    description: 'Tracetest attribute that reflects the elapsed real time of the operation.',
    note: '',
    tags: ['ms', 'second', 'time'],
  },
  'tracetest.span.type': {
    description:
      'Tracetest attribute based on the [OTel Trace Semantic Conventions](https://github.com/open-telemetry/semantic-conventions/blob/main/docs/README.md)',
    note: '',
    tags: ['general', 'http', 'database', 'rpc', 'messaging', 'faas', 'exception'],
  },
};
