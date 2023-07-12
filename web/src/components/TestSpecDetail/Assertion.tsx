import AttributeValue from 'components/AttributeValue';
import AssertionService from 'services/Assertion.service';
import OperatorService from 'services/Operator.service';
import {ICheckResult} from 'types/Assertion.types';
import {TCompareOperatorSymbol} from 'types/Operator.types';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from '../Editor';
import * as S from './TestSpecDetail.styled';

interface IProps {
  check: ICheckResult;
  testId: string;
  runId: string;
  selector: string;
}

const Assertion = ({check, testId, runId, selector}: IProps) => (
  <S.CheckItemContainer>
    <S.GridContainer>
      {check.result.error && AssertionService.isValidError(check.result.error) ? (
        <>
          <S.Row $justify="center">
            <S.IconWarning />
          </S.Row>
          <AttributeValue strong type="warning" value={check.result.error} />
        </>
      ) : (
        <>
          <S.Row $justify="center">{check.result.passed ? <S.IconSuccess /> : <S.IconError />}</S.Row>
          <AttributeValue
            strong
            type={check.result.passed ? 'success' : 'danger'}
            value={check.result.observedValue || '<Empty Value>'}
          />
        </>
      )}
      <div />
      <S.Row>
        <S.AssertionContainer>
          <Editor
            type={SupportedEditors.Expression}
            value={check.assertion}
            editable={false}
            context={{
              testId,
              runId,
              spanId: check.result.spanId,
              selector,
            }}
          />
          <S.SecondaryText>
            {OperatorService.getNameFromSymbol(check.assertion as TCompareOperatorSymbol)}
          </S.SecondaryText>
        </S.AssertionContainer>
      </S.Row>
    </S.GridContainer>
  </S.CheckItemContainer>
);

export default Assertion;
