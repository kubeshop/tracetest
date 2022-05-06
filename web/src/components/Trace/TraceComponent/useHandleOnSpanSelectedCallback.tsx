import {Dispatch, SetStateAction, useCallback} from 'react';
import {ActionCreator} from '@reduxjs/toolkit';
import {Elements} from 'react-flow-renderer';
import {ISpan} from 'types/Span.types';
import {ITestRunResult} from 'types/TestRunResult.types';

export function useHandleOnSpanSelectedCallback(
  addSelected: ActionCreator<{type: 'ADD_SELECTED_ELEMENTS'; payload: Elements<any>}>,
  testResultDetails: ITestRunResult | undefined,
  setSelectedSpan: Dispatch<SetStateAction<ISpan | undefined>>
) {
  return useCallback(
    (spanId: string) => {
      addSelected([{id: spanId}]);
      setSelectedSpan(testResultDetails?.trace?.spans.find(({spanId: id}) => id === spanId));
    },
    [addSelected, setSelectedSpan, testResultDetails?.trace?.spans]
  );
}
