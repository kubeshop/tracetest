import useHover from 'hooks/useHover';
import {useEffect} from 'react';
import {IResult} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import AttributeValue from '../AttributeValue';
import {Steps} from '../GuidedTour/traceStepList';
import AttributeCheck from './AttributeCheck';
import * as S from './AttributeRow.styled';

interface IProps {
  assertionsFailed?: IResult[];
  assertionsPassed?: IResult[];
  attribute: TSpanFlatAttribute;
  isCopied: boolean;
  onCopy(value: string): void;
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
  setIsCopied(value: boolean): void;
  shouldDisplayActions: boolean;
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
  shouldDisplayActions,
}: IProps) => {
  const {isHovering, onMouseEnter, onMouseLeave} = useHover();
  const passedCount = assertionsPassed?.length ?? 0;
  const failedCount = assertionsFailed?.length ?? 0;

  useEffect(() => {
    if (!isHovering) setIsCopied(false);
  }, [isHovering, setIsCopied]);

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
        {(isHovering || shouldDisplayActions) && (
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
