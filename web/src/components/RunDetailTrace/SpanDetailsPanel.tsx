import {useNavigate} from 'react-router-dom';
import {useCallback} from 'react';
import {useAppSelector} from 'redux/hooks';
import TraceSelectors from 'selectors/Trace.selectors';
import TestRun from 'models/TestRun.model';
import SpanSelectors from 'selectors/Span.selectors';
import {TPanel, TPanelComponentProps} from '../ResizablePanels/ResizablePanels';
import SpanDetail from '../SpanDetail/SpanDetail';
import * as S from './RunDetailTrace.styled';

type TProps = TPanelComponentProps & {
  run: TestRun;
  testId: string;
};

const SpanDetailsPanel = ({size: {isOpen}, run, testId}: TProps) => {
  const searchText = useAppSelector(TraceSelectors.selectSearchText);
  const selectedSpan = useAppSelector(TraceSelectors.selectSelectedSpan);
  const navigate = useNavigate();
  const span = useAppSelector(state => SpanSelectors.selectSpanById(state, selectedSpan, testId, run.id));

  const handleOnCreateSpec = useCallback(() => {
    navigate(`/test/${testId}/run/${run.id}/test`);
  }, [navigate, run.id, testId]);

  return (
    <S.PanelContainer $isOpen={isOpen}>
      <SpanDetail onCreateTestSpec={handleOnCreateSpec} searchText={searchText} span={span} />
    </S.PanelContainer>
  );
};

export const getSpanDetailsPanel = (testId: string, run: TestRun, order = 1): TPanel => ({
  name: `SPAN_DETAILS_${order}`,
  minSize: 15,
  maxSize: 320,
  position: 'left',
  component: props => <SpanDetailsPanel {...props} testId={testId} run={run} />,
});

export default SpanDetailsPanel;
