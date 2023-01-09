import {Button} from 'antd';

import TestOutput from 'components/TestOutput';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import Empty from './Empty';
import * as S from './TestOutputs.styled';
import {TTestOutput} from '../../types/TestOutput.types';

interface IProps {
  outputs: TTestOutput[];
}

const TestOutputs = ({outputs}: IProps) => {
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
