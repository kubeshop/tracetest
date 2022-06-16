import {useCallback} from 'react';
import {useStoreActions} from 'react-flow-renderer';

import RunBottomPanel from 'components/RunBottomPanel';
import {RunLayoutProvider} from 'components/RunLayout';
import RunTopPanel from 'components/RunTopPanel';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setSelectedSpan} from 'redux/slices/TestDefinition.slice';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TTest} from 'types/Test.types';
import {TTestRun} from 'types/TestRun.types';

interface IProps {
  displayError: boolean;
  run: TTestRun;
  test?: TTest;
}

const Run = ({displayError, run, test}: IProps) => {
  const dispatch = useAppDispatch();
  const addSelected = useStoreActions(actions => actions.addSelectedElements);
  const selectedSpan = useAppSelector(TestDefinitionSelectors.selectSelectedSpan);

  const onSelectSpan = useCallback(
    (spanId: string) => {
      const span = run?.trace?.spans.find(({id}) => id === spanId);
      if (span) addSelected([{id: span?.id}]);
      dispatch(setSelectedSpan(span));
    },
    [addSelected, dispatch, run?.trace?.spans]
  );

  if (displayError) {
    return null;
  }

  return (
    <RunLayoutProvider
      bottomPanel={
        <RunBottomPanel onSelectSpan={onSelectSpan} run={run} selectedSpan={selectedSpan!} testId={test?.id!} />
      }
      topPanel={<RunTopPanel onSelectSpan={onSelectSpan} run={run} selectedSpan={selectedSpan!} />}
    />
  );
};

export default Run;
