import {Button} from 'antd';

import TestOutput from 'components/TestOutput';
import {useTest} from 'providers/Test/Test.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import {selectTestOutputs} from 'redux/testOutputs/selectors';
import Empty from './Empty';
import * as S from './TestOutputs.styled';

const TestOutputs = () => {
  const {
    run: {id: runId},
  } = useTestRun();
  const {
    test: {id: testId},
  } = useTest();
  const outputs = useAppSelector(state => selectTestOutputs(state, testId, runId));
  const {onDelete, onOpen} = useTestOutput();

  return (
    <S.Container>
      <S.HeaderContainer>
        <Button data-cy="output-add-button" onClick={() => onOpen()} type="primary">
          Add Test Output
        </Button>
      </S.HeaderContainer>

      {!outputs.length && <Empty />}

      <S.ListContainer>
        {outputs.map((output, index) => (
          <TestOutput index={index} key={output.name} output={output} onDelete={onDelete} onEdit={onOpen} />
        ))}
      </S.ListContainer>
    </S.Container>
  );
};

export default TestOutputs;
