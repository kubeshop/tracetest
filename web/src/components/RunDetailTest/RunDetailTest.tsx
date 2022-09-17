import Drawer from 'components/Drawer';
import SpanDetail from 'components/SpanDetail';
import TestResults from 'components/TestResults';
import TestSpecForm from 'components/TestSpecForm';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {useSpan} from 'providers/Span/Span.provider';
import {TTestRun} from 'types/TestRun.types';
import { useTestSpecs } from '../../providers/TestSpecs/TestSpecs.provider';
import * as S from './RunDetailTest.styled';
import Visualization from './Visualization';

interface IProps {
  run: TTestRun;
  testId: string;
}

const RunDetailTest = ({run, testId}: IProps) => {
  const {selectedSpan} = useSpan();
  const {selectedTestSpec} = useTestSpecs();
  const {isOpen: isTestSpecFormOpen, formProps, onSubmit, close} = useTestSpecForm();

  return (
    <S.Container>
      <Drawer>
        <SpanDetail span={selectedSpan} />
      </Drawer>

      <S.Container>
        <S.SectionLeft>
          <Visualization runState={run.state} spans={run?.trace?.spans ?? []} />
        </S.SectionLeft>
        <S.SectionRight $shouldScroll={!selectedTestSpec}>
          {isTestSpecFormOpen ? (
            <TestSpecForm
              onSubmit={onSubmit}
              runId={run.id}
              testId={testId}
              {...formProps}
              onCancel={() => {
                close();
              }}
            />
          ) : (
            <TestResults />
          )}
        </S.SectionRight>
      </S.Container>
    </S.Container>
  );
};

export default RunDetailTest;
