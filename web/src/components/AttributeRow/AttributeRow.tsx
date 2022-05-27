import AttributeValue from 'components/AttributeValue';
import useHover from 'hooks/useHover';
import {TSpanFlatAttribute} from 'types/Span.types';
import {IResult} from 'components/SpanDetail/SpanDetail';
import AttributeCheck from './AttributeCheck';
import * as S from './AttributeRow.styled';

interface IProps {
  assertionsFailed?: IResult[];
  assertionsPassed?: IResult[];
  attribute: TSpanFlatAttribute;
  onCopy(value: string): void;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
}

const AttributeRow = ({
  assertionsFailed,
  assertionsPassed,
  attribute: {key, value},
  attribute,
  onCopy,
  onCreateAssertion,
}: IProps) => {
  const {isHovering, onMouseEnter, onMouseLeave} = useHover();
  const passedCount = assertionsPassed?.length ?? 0;
  const failedCount = assertionsFailed?.length ?? 0;

  return (
    <S.AttributeRow onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
      <S.TextContainer>
        <S.Text type="secondary">{key}</S.Text>
      </S.TextContainer>

      <S.AttributeValueRow>
        <AttributeValue value={value} />
        {passedCount > 0 && <AttributeCheck items={assertionsPassed!} type="success" />}
        {failedCount > 0 && <AttributeCheck items={assertionsFailed!} type="error" />}
      </S.AttributeValueRow>

      <S.IconContainer>
        {isHovering && (
          <>
            <S.CopyIcon onClick={() => onCopy(value)} />
            <S.CustomTooltip placement="top" title="Add Assertion">
              <S.AddAssertionIcon onClick={() => onCreateAssertion(attribute)} />
            </S.CustomTooltip>
          </>
        )}
      </S.IconContainer>
    </S.AttributeRow>
  );
};

export default AttributeRow;
