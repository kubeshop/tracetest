import FlamegraphService from 'services/Flamegraph.service';
import {TSpan} from 'types/Span.types';

describe('Flame View Test', () => {
  const traceId = '1';
  const s1 = {id: '1', parentId: '', name: 'one'} as TSpan;
  const s2 = {id: '2', parentId: '1', name: 'two'} as TSpan;
  const s3 = {id: '3', parentId: '2', name: 'three'} as TSpan;
  const spans = [s1, s2, s3];

  test('parentReferences', () => {
    expect(FlamegraphService.parentReferences(traceId, spans, s2)).toStrictEqual([
      {
        refType: 'CHILD_OF',
        spanID: s1.id,
        traceID: traceId,
      },
    ]);
  });

  test('tracetestSpanToJagerSpan', () => {
    expect(FlamegraphService.tracetestSpanToJagerSpan(traceId, spans, s1)).toStrictEqual({
      operationName: s1.name,
      duration: 0,
      processID: '',
      references: FlamegraphService.parentReferences(traceId, spans, s1),
      spanID: s1.id,
      startTime: 0,
      tags: [],
      traceID: traceId,
      warnings: [],
      flags: '',
      logs: {fields: [], timestamp: 0},
    });
  });
});
