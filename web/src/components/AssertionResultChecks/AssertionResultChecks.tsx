import {IResult} from 'types/Assertion.types';
import AttributeCheck from './ResultCheck';
import * as S from './AssertionResultChecks.styled';

interface IProps {
  failed: IResult[];
  passed: IResult[];
  styleType?: 'node' | 'summary' | 'default';
}

const AssertionResultChecks = ({passed, failed, styleType}: IProps) => (
  <S.Container>
    {Boolean(passed.length) && <AttributeCheck items={passed} type="success" styleType={styleType} />}
    {Boolean(failed.length) && <AttributeCheck items={failed} type="error" styleType={styleType} />}
  </S.Container>
);

export default AssertionResultChecks;
