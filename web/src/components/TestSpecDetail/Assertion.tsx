import AttributeValue from 'components/AttributeValue';
import OperatorService from 'services/Operator.service';
import {ICheckResult} from 'types/Assertion.types';
import {TCompareOperatorSymbol} from 'types/Operator.types';
import {SupportedEditors} from 'constants/Editor.constants';
import Editor from '../Editor';
import * as S from './TestSpecDetail.styled';

interface IProps {
  check: ICheckResult;
}

const CheckItem = ({check}: IProps) => (
  <S.CheckItemContainer>
    <S.GridContainer>
      <S.Row>{check.result.passed ? <S.IconSuccess /> : <S.IconError />}</S.Row>
      <S.Row>
        <S.AssertionContainer>
          <Editor type={SupportedEditors.Expression} value={check.assertion} editable={false} />
          <S.SecondaryText>
            {OperatorService.getNameFromSymbol(check.assertion as TCompareOperatorSymbol)}
          </S.SecondaryText>
        </S.AssertionContainer>
      </S.Row>
    </S.GridContainer>
    <S.GridContainer>
      <div />
      <AttributeValue
        strong
        type={check.result.passed ? 'success' : 'danger'}
        value={check.result.observedValue || '<Empty Value>'}
      />
    </S.GridContainer>
  </S.CheckItemContainer>
);

export default CheckItem;
