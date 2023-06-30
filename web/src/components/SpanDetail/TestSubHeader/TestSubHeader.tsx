import AssertionResultChecks from 'components/AssertionResultChecks';
import TestOutputMark from 'components/TestOutputMark';
import TestOutput from 'models/TestOutput.model';
import * as S from './TestSubHeader.styled';
import {IPropsSubHeader} from '../SpanDetail';

const TestSubHeader = ({testSpecs, testOutputs}: IPropsSubHeader) => {
  if (!testSpecs && !testOutputs?.length) return null;

  return (
    <S.Container>
      {testSpecs && <AssertionResultChecks failed={testSpecs.failed} passed={testSpecs.passed} styleType="summary" />}
      {!!testOutputs?.length && <TestOutputMark outputs={testOutputs as TestOutput[]} />}
    </S.Container>
  );
};

export default TestSubHeader;
