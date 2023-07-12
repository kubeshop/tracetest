import {TTestSpec} from 'types/TestRun.types';
import ResultCheck from './ResultCheck';
import * as S from './AssertionResultChecks.styled';

interface IProps {
  failed: TTestSpec[];
  passed: TTestSpec[];
  styleType?: 'node' | 'summary' | 'default';
}

const AssertionResultChecks = ({passed, failed, styleType}: IProps) => (
  <S.Container>
    {Boolean(passed.length) && <ResultCheck items={passed} type="success" styleType={styleType} />}
    {Boolean(failed.length) && <ResultCheck items={failed} type="error" styleType={styleType} />}
  </S.Container>
);

export default AssertionResultChecks;
