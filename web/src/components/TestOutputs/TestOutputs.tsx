import {Button} from 'antd';

import TestOutput from 'components/TestOutput';
import TestOutputModel from 'models/TestOutput.model';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import SpanService from 'services/Span.service';
import Empty from './Empty';
import * as S from './TestOutputs.styled';

interface IProps {
  outputs: TestOutputModel[];
}

const TestOutputs = ({outputs}: IProps) => {
  const {onDelete, onOpen} = useTestOutput();
  const {selectedSpan} = useSpan();

  const handleOnAddTestOutput = () => {
    const query = selectedSpan ? SpanService.getSelectorInformation(selectedSpan) : '';
    onOpen(TestOutputModel({selector: query, selectorParsed: {query}}));
  };

  return (
    <S.Container>
      <S.HeaderContainer>
        <Button data-cy="output-add-button" onClick={handleOnAddTestOutput} type="primary">
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
