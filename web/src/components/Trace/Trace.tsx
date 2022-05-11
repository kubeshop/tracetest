import {useState} from 'react';
import {Tabs} from 'antd';
import {useStoreActions} from 'react-flow-renderer';
import {ISpan} from 'types/Span.types';
import {ITestRunResult} from 'types/TestRunResult.types';
import {ITest} from 'types/Test.types';
import Diagram, {SupportedDiagrams} from 'components/Diagram/Diagram';
import SpanDetail from 'components/SpanDetail';
import {TimelineDrawer} from './TimelineDrawer';
import {useHandleOnSpanSelectedCallback} from './hooks/useHandleOnSpanSelectedCallback';
import * as S from './Trace.styled';
import TestResults from './TestResults';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {Steps} from '../GuidedTour/traceStepList';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';
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
          <S.TabsContainer>
            <S.TraceTabs onChange={activeTab => TraceAnalyticsService.onChangeTab(activeTab)}>
              <Tabs.TabPane
                tab={
                  <span data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.SpanDetail)}>Span Detail</span>
                }
                key="span-detail"
              >
                <SpanDetail resultId={testResultDetails?.resultId} testId={test?.testId} span={selectedSpan} />
              </Tabs.TabPane>
              <Tabs.TabPane
                tab={
                  <span data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.TestResults)}>Test Results</span>
                }
                key="test-results"
              >
                <TestResults
                  onSpanSelected={onSelectSpan}
                  trace={testResultDetails?.trace}
                  resultId={testResultDetails?.resultId!}
                />
              </Tabs.TabPane>
            </S.TraceTabs>
          </S.TabsContainer>
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
