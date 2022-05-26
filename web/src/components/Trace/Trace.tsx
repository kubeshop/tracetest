import {useState} from 'react';
import {useStoreActions} from 'react-flow-renderer';
import {TSpan} from 'types/Span.types';
import {TTestRun} from 'types/TestRun.types';
import {TTest} from 'types/Test.types';
import Diagram, {SupportedDiagrams} from 'components/Diagram/Diagram';
import SpanDetail from 'components/SpanDetail';
import {useHandleOnSpanSelectedCallback} from './hooks/useHandleOnSpanSelectedCallback';
import * as S from './Trace.styled';
import DiagramSwitcher from '../DiagramSwitcher';
import TraceDrawer from '../TraceDrawer';

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
          <Diagram type={diagramType} trace={run.trace!} onSelectSpan={onSelectSpan} selectedSpan={selectedSpan} />
        </S.DiagramSection>
        <S.DetailsSection>
          <SpanDetail resultId={run.id} testId={test?.id} span={selectedSpan} />
        </S.DetailsSection>
      </S.Main>
      <TraceDrawer selectedSpan={selectedSpan!} visiblePortion={visiblePortion} testId={test?.id!} run={run} onSelectSpan={onSelectSpan} />
    </>
  ) : null;
};

export default Trace;
