import {useEffect} from 'react';

import useHover from 'hooks/useHover';
import {IResult} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import AttributeCheck from './AttributeCheck';
import * as S from './AttributeRow.styled';
import AttributeValue from '../AttributeValue';

interface IProps {
  assertionsFailed?: IResult[];
  assertionsPassed?: IResult[];
  attribute: TSpanFlatAttribute;
  isCopied: boolean;
  onCopy(value: string): void;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
  setIsCopied(value: boolean): void;
}

const AttributeRow = ({
  assertionsFailed,
  assertionsPassed,
  attribute: {key, value},
  attribute,
  isCopied,
  onCopy,
  onCreateAssertion,
  setIsCopied,
}: IProps) => {
  const {isHovering, onMouseEnter, onMouseLeave} = useHover();
  const passedCount = assertionsPassed?.length ?? 0;
  const failedCount = assertionsFailed?.length ?? 0;

  useEffect(() => {
    if (!isHovering) setIsCopied(false);
  }, [isHovering]);

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
            <S.CustomTooltip title={isCopied ? 'Copied' : 'Copy'}>
              <S.CopyIcon onClick={() => onCopy(value)} />
            </S.CustomTooltip>
            <S.CustomTooltip title="Add Assertion">
              <S.AddAssertionIcon onClick={() => onCreateAssertion(attribute)} />
            </S.CustomTooltip>
          </>
        )}
      </S.IconContainer>
    </S.AttributeRow>
  );
};

export default AttributeRow;
