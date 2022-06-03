import {useState} from 'react';
import {useStoreActions} from 'react-flow-renderer';

import Diagram, {SupportedDiagrams} from 'components/Diagram/Diagram';
import DiagramSwitcher from 'components/DiagramSwitcher';
import SpanDetail from 'components/SpanDetail';
import TraceDrawer from 'components/TraceDrawer';
import {useAppSelector} from 'redux/hooks';
import TestDefinitionSelectors from 'selectors/TestDefinition.selectors';
import {TSpan} from 'types/Span.types';
import {TTest} from 'types/Test.types';
import {TTestRun} from 'types/TestRun.types';
import {useHandleOnSpanSelectedCallback} from './hooks/useHandleOnSpanSelectedCallback';
import * as S from './Trace.styled';

interface IProps {
  displayError: boolean;
  minHeight: string;
  run: TTestRun;
  test?: TTest;
  visiblePortion: number;
}

const Trace = ({displayError, visiblePortion, minHeight, test, run}: IProps): JSX.Element | null => {
  const [selectedSpan, setSelectedSpan] = useState<TSpan | undefined>();
  const [diagramType, setDiagramType] = useState<SupportedDiagrams>(SupportedDiagrams.DAG);
  const affectedSpans = useAppSelector(TestDefinitionSelectors.selectAffectedSpans);

  const addSelected = useStoreActions(actions => actions.addSelectedElements);
  const onSelectSpan = useHandleOnSpanSelectedCallback(addSelected, run, setSelectedSpan);

  return !displayError ? (
    <>
      <S.Main height={minHeight}>
        <S.DiagramSection>
          <DiagramSwitcher
            onTypeChange={setDiagramType}
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
