import {DoubleLeftOutlined, DoubleRightOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import {useState} from 'react';

import Diagram from 'components/Diagram';
import {SupportedDiagrams} from 'components/Diagram/Diagram';
import SpanDetail from 'components/SpanDetail';
import TestResults from 'components/TestResults';
import TestSpecForm from 'components/TestSpecForm';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {useSpan} from 'providers/Span/Span.provider';
import {TTestRun} from 'types/TestRun.types';
import * as S from './RunDetailTest.styled';

interface IProps {
  run: TTestRun;
  testId: string;
}

const RunDetailTest = ({run, testId}: IProps) => {
  const {selectedSpan} = useSpan();
  const {isOpen: isTestSpecFormOpen, formProps, onSubmit, close} = useTestSpecForm();
  const [isAsideOpen, setIsAsideOpen] = useState(false);

  return (
    <S.Container>
      <S.Aside $isOpen={isAsideOpen}>
        <S.AsideContent>
          <SpanDetail span={selectedSpan} />
        </S.AsideContent>
        <S.AsideButtonContainer>
          <Button
            icon={isAsideOpen ? <DoubleLeftOutlined /> : <DoubleRightOutlined />}
            onClick={() => setIsAsideOpen(isOpen => !isOpen)}
            shape="circle"
            size="small"
            type="primary"
          />
        </S.AsideButtonContainer>
      </S.Aside>

      <S.Container>
        <S.SectionLeft>
          <Diagram trace={run.trace!} runState={run.state} type={SupportedDiagrams.DAG} />
        </S.SectionLeft>
        <S.SectionRight>
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
