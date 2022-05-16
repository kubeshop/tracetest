import {useCallback} from 'react';
import useHover from '../../hooks/useHover';
import {ISpanFlatAttribute} from '../../types/Span.types';
import * as S from './AttributeRow.styled';
import AttributeValue from '../AttributeValue';

interface IAttributeRowProps {
  attribute: ISpanFlatAttribute;
  onCreateAssertion(attribute: ISpanFlatAttribute): void;
}

const AttributeRow: React.FC<IAttributeRowProps> = ({attribute: {key, value}, attribute, onCreateAssertion}) => {
  const {isHovering, onMouseEnter, onMouseLeave} = useHover();

  const onCopy = useCallback(() => {
    navigator.clipboard.writeText(value);
  }, [value]);

  return (
    <S.AttributeRow onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
      <S.TextContainer>
        <S.Text type="secondary">{key}</S.Text>
      </S.TextContainer>
      <AttributeValue value={value} />
      <S.IconContainer>
        {isHovering && (
          <>
            <S.CopyIcon onClick={onCopy} />
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
