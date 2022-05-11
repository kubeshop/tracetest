import {useState} from 'react';
import {useStoreActions} from 'react-flow-renderer';
import {ISpan} from 'types/Span.types';
import {ITestRunResult} from 'types/TestRunResult.types';
import {ITest} from 'types/Test.types';
import Diagram, {SupportedDiagrams} from 'components/Diagram/Diagram';
import SpanDetail from 'components/SpanDetail';
import {TimelineDrawer} from './TimelineDrawer';
import {useHandleOnSpanSelectedCallback} from './hooks/useHandleOnSpanSelectedCallback';
import * as S from './Trace.styled';
import DiagramSwitcher from '../DiagramSwitcher';

interface IProps {
  displayError: boolean;
  minHeight: string;
  testResultDetails: ITestRunResult | undefined;
  test?: ITest;
  visiblePortion: number;
}

const Trace = ({displayError, visiblePortion, minHeight, test, testResultDetails}: IProps): JSX.Element | null => {
  const [selectedSpan, setSelectedSpan] = useState<ISpan | undefined>();
  const [diagramType, setDiagramType] = useState<SupportedDiagrams>(SupportedDiagrams.DAG);

  const addSelected = useStoreActions(actions => actions.addSelectedElements);
  const onSelectSpan = useHandleOnSpanSelectedCallback(addSelected, testResultDetails, setSelectedSpan);

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
            type={diagramType}
            trace={testResultDetails?.trace!}
            onSelectSpan={onSelectSpan}
            selectedSpan={selectedSpan}
          />
        </S.DiagramSection>
        <S.DetailsSection>
          <SpanDetail resultId={testResultDetails?.resultId} testId={test?.testId} span={selectedSpan} />
        </S.DetailsSection>
      </S.Main>
      <TimelineDrawer
        visiblePortion={visiblePortion}
        testResultDetails={testResultDetails}
        onSelectSpan={onSelectSpan}
        selectedSpan={selectedSpan}
      />
    </>
  ) : null;
};

export default Trace;
