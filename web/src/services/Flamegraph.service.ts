import {Trace} from '@pyroscope/models';
import {TSpan} from '../types/Span.types';

type Span = Trace['spans'][number];

const FlameGraphService = () => ({
  parentReferences(traceId: string, spans: TSpan[], current: TSpan): Span['references'] {
    const parent = spans.find(s => s.id === current.parentId);
    if (!parent) return [];
    return [
      {
        refType: 'CHILD_OF',
        spanID: parent.id,
        traceID: traceId,
      },
    ];
  },
  tracetestSpanToJagerSpan(id: string, spans: TSpan[], s: TSpan): Span {
    const milliseconds = 1000;
    return {
      operationName: s.name,
      duration: ((s.endTime || 0) - (s.startTime || 0)) * milliseconds,
      processID: '',
      references: this.parentReferences(id, spans, s),
      spanID: s.id,
      startTime: (s.startTime || 0) * milliseconds,
      tags: [],
      traceID: id,
      warnings: [],
      flags: '',
      logs: {fields: [], timestamp: 0},
    };
  },
  convertTracetestSpanToJaeger(traceID: string, spans: TSpan[]): Trace {
    return {
      traceID,
      spans: spans.map(s => this.tracetestSpanToJagerSpan(traceID, spans, s)),
      processes: {p1: {serviceName: '', tags: []}, p2: {serviceName: '', tags: []}},
    };
  },
});

export default FlameGraphService();
