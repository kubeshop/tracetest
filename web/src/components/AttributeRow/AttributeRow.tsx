import useHover from 'hooks/useHover';
import {Dispatch, SetStateAction, useEffect} from 'react';
import {IResult} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import AttributeValue from '../AttributeValue';
import {Steps} from '../GuidedTour/traceStepList';
import AttributeCheck from './AttributeCheck';
import * as S from './AttributeRow.styled';
import {useHoverComponent} from './useHoverComponent';

interface IProps {
  isAnyHovered: number[];
  setIsAnyHovered: Dispatch<SetStateAction<number[]>>;
  assertionsFailed?: IResult[];
  index: number;
  assertionsPassed?: IResult[];
  attribute: TSpanFlatAttribute;
  isCopied: boolean;
  onCopy(value: string): void;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
  setIsCopied(value: boolean): void;
}

const AttributeRow = ({
  isAnyHovered,
  setIsAnyHovered,
  assertionsFailed,
  assertionsPassed,
  attribute: {key, value},
  attribute,
  isCopied,
  onCopy,
  onCreateAssertion,
  setIsCopied,
  index,
}: IProps) => {
  const isFirst = index === 0;

  const {isHovering, onMouseEnter, onMouseLeave} = useHover();
  const isAnyHoverCallback = useHoverComponent(index, isAnyHovered, setIsAnyHovered);
  const passedCount = assertionsPassed?.length ?? 0;
  const failedCount = assertionsFailed?.length ?? 0;

  useEffect(() => {
    if (!isHovering) setIsCopied(false);
  }, [isHovering, setIsCopied]);

  return (
    <S.AttributeRow
      onMouseEnter={() => {
        onMouseEnter();
        isAnyHoverCallback.onMouseEnter();
      }}
      onMouseLeave={() => {
        onMouseLeave();
        isAnyHoverCallback.onMouseLeave();
      }}
    >
      <S.TextContainer>
        <S.Text type="secondary">{key}</S.Text>
      </S.TextContainer>

      <S.AttributeValueRow>
        <AttributeValue value={value} />
        {passedCount > 0 && <AttributeCheck items={assertionsPassed!} type="success" />}
        {failedCount > 0 && <AttributeCheck items={assertionsFailed!} type="error" />}
      </S.AttributeValueRow>

      <S.IconContainer>
        {(isHovering || (isFirst && isAnyHovered.length === 0)) && (
          <>
            <S.CustomTooltip title={isCopied ? 'Copied' : 'Copy'}>
              <S.CopyIcon onClick={() => onCopy(value)} />
            </S.CustomTooltip>
            <S.CustomTooltip title="Add Assertion">
              <S.AddAssertionIcon
                data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Assertions)}
                onClick={() => onCreateAssertion(attribute)}
              />
            </S.CustomTooltip>
          </>
        )}
      </S.IconContainer>
    </S.AttributeRow>
  );
};

export default AttributeRow;
