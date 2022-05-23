import {Dispatch, SetStateAction, useCallback} from 'react';
import {ActionCreator} from '@reduxjs/toolkit';
import {Elements} from 'react-flow-renderer';
import {TSpan} from 'types/Span.types';
import {TTestRun} from 'types/TestRun.types';

export function useHandleOnSpanSelectedCallback(
  addSelected: ActionCreator<{type: 'ADD_SELECTED_ELEMENTS'; payload: Elements<any>}>,
  testResultDetails: TTestRun | undefined,
  setSelectedSpan: Dispatch<SetStateAction<TSpan | undefined>>
) {
  return useCallback(
    (spanId: string) => {
      addSelected([{id: spanId}]);
      setSelectedSpan(testResultDetails?.trace?.spans.find(({id}) => id === spanId));
    },
    [addSelected, setSelectedSpan, testResultDetails?.trace?.spans]
  );
}
