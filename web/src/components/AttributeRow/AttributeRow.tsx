import {useCallback, useEffect, useState} from 'react';
import useHover from '../../hooks/useHover';
import {TSpanFlatAttribute} from '../../types/Span.types';
import * as S from './AttributeRow.styled';
import AttributeValue from '../AttributeValue';

interface IAttributeRowProps {
  attribute: TSpanFlatAttribute;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
}

const AttributeRow: React.FC<IAttributeRowProps> = ({attribute: {key, value}, attribute, onCreateAssertion}) => {
  const {isHovering, onMouseEnter, onMouseLeave} = useHover();
  const [isCopy, setIsCopy] = useState(false);

  const onCopy = useCallback(() => {
    navigator.clipboard.writeText(value);
    setIsCopy(true);
  }, [value]);

  useEffect(() => {
    if (!isHovering) setIsCopy(false);
  }, [isHovering]);

  return (
    <S.AttributeRow onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
      <S.TextContainer>
        <S.Text type="secondary">{key}</S.Text>
      </S.TextContainer>
      <AttributeValue value={value} />
      <S.IconContainer>
        {isHovering && (
          <>
            <S.CustomTooltip title={isCopy ? 'Copied' : 'Copy'}>
              <S.CopyIcon onClick={onCopy} />
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
