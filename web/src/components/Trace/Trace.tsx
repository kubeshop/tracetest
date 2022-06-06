import {useCallback, useState} from 'react';
import {useStoreActions} from 'react-flow-renderer';

import Diagram, {SupportedDiagrams} from 'components/Diagram/Diagram';
import DiagramSwitcher from 'components/DiagramSwitcher';
import SpanDetail from 'components/SpanDetail';
import TraceDrawer from 'components/TraceDrawer';
import {useAppDispatch, useAppSelector} from 'redux/hooks';
import {setSelectedSpan} from 'redux/slices/TestDefinition.slice';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TTest} from 'types/Test.types';
import {TTestRun} from 'types/TestRun.types';
import * as S from './Trace.styled';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';

interface IProps {
  displayError: boolean;
  minHeight: string;
  run: TTestRun;
  test?: TTest;
  visiblePortion: number;
}

const Trace = ({displayError, visiblePortion, minHeight, test, run}: IProps): JSX.Element | null => {
  const dispatch = useAppDispatch();
  const addSelected = useStoreActions(actions => actions.addSelectedElements);
  const [diagramType, setDiagramType] = useState<SupportedDiagrams>(SupportedDiagrams.DAG);
  const selectedSpan = useAppSelector(TestDefinitionSelectors.selectSelectedSpan);
  const affectedSpans = useAppSelector(TestDefinitionSelectors.selectAffectedSpans);

  const onSelectSpan = useCallback(
    (spanId: string) => {
      const span = run?.trace?.spans.find(({id}) => id === spanId);
      if (span) addSelected([{id: span?.id}]);
      dispatch(setSelectedSpan(span));
    },
    [addSelected, dispatch, run?.trace?.spans]
  );

  return !displayError ? (
    <>
      <S.Main height={minHeight}>
        <S.DiagramSection>
          <DiagramSwitcher
            onTypeChange={type => {
              TraceAnalyticsService.onSwitchDiagramView(type);
              setDiagramType(type);
            }}
            onSearch={() => console.log('onSearch')}
            selectedType={diagramType}
          />
          <Diagram
            affectedSpans={affectedSpans}
            onSelectSpan={onSelectSpan}
            selectedSpan={selectedSpan}
            trace={run.trace!}
            runState={run.state}
            type={diagramType}
          />
        </S.DiagramSection>
        <S.DetailsSection>
          <SpanDetail span={selectedSpan} />
        </S.DetailsSection>
      </S.Main>
      <TraceDrawer
        selectedSpan={selectedSpan!}
        visiblePortion={visiblePortion}
        testId={test?.id!}
        run={run}
        onSelectSpan={onSelectSpan}
      />
    </>
  ) : null;
};

export default Trace;
