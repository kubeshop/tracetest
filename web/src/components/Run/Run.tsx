import {useCallback} from 'react';
import {useStoreActions} from 'react-flow-renderer';

import RunBottomPanel from 'components/RunBottomPanel';
import {RunLayoutProvider} from 'components/RunLayout';
import RunTopPanel from 'components/RunTopPanel';
import {TTest} from 'types/Test.types';
import {TTestRun} from 'types/TestRun.types';
import {useSpan} from '../../providers/Span/Span.provider';

interface IProps {
  displayError: boolean;
  run: TTestRun;
  test?: TTest;
}

const Run = ({displayError, run, test}: IProps) => {
  const addSelected = useStoreActions(actions => actions.addSelectedElements);
  const {selectedSpan, onSelectSpan} = useSpan();

  const handleSelectSpan = useCallback(
    (spanId: string) => {
      const span = run?.trace?.spans.find(({id}) => id === spanId);
      if (span) {
        addSelected([{id: span?.id}]);
        onSelectSpan(span);
      }
    },
    [addSelected, onSelectSpan, run?.trace?.spans]
  );

  if (displayError) {
    return null;
  }

  return (
    <RunLayoutProvider
      bottomPanel={
        <RunBottomPanel onSelectSpan={handleSelectSpan} run={run} selectedSpan={selectedSpan!} testId={test?.id!} />
      }
      topPanel={<RunTopPanel onSelectSpan={handleSelectSpan} run={run} selectedSpan={selectedSpan!} />}
    />
  );
};

export default Run;
