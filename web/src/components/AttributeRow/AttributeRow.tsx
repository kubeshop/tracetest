import {useCallback, useMemo, useState} from 'react';
import JSONPretty from 'react-json-pretty';
import useHover from '../../hooks/useHover';
import {ISpanFlatAttribute} from '../../types/Span.types';
import {isJson} from '../../utils/Common';
import * as S from './AttributeRow.styled';

interface IAttributeRowProps {
  attribute: ISpanFlatAttribute;
  onCreateAssertion(attribute: ISpanFlatAttribute): void;
}

const AttributeRow: React.FC<IAttributeRowProps> = ({attribute: {key, value}, attribute, onCreateAssertion}) => {
  const {isHovering, onMouseEnter, onMouseLeave} = useHover();
  const [isCollapsed, setIsCollapsed] = useState(false);

  const onCopy = useCallback(() => {
    navigator.clipboard.writeText(value);
  }, [value]);

  const parsedValue = useMemo(() => (isJson(value) ? <JSONPretty data={value} /> : value), [value]);

  return (
    <S.AttributeRow onMouseEnter={onMouseEnter} onMouseLeave={onMouseLeave}>
      <S.TextContainer>
        <S.Text type="secondary">{key}</S.Text>
      </S.TextContainer>
      <S.ValueText onClick={() => setIsCollapsed(!isCollapsed)} isCollapsed={isCollapsed}>
        {parsedValue}
      </S.ValueText>
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
