import {Steps} from 'components/GuidedTour/traceStepList';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import SpanService from 'services/Span.service';
import {TSpan} from 'types/Span.types';
import {singularOrPlural} from 'utils/Common';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import * as S from './TestResults.styled';

interface IProps {
  selectedSpan: TSpan;
  totalFailedSpecs: number;
  totalPassedSpecs: number;
}

const Header = ({selectedSpan, totalFailedSpecs, totalPassedSpecs}: IProps) => {
  const {open} = useTestSpecForm();

  const handleAddTestSpecOnClick = () => {
    const selector = SpanService.getSelectorInformation(selectedSpan!);
    open({
      isEditing: false,
      selector,
      defaultValues: {
        selector,
      },
    });
  };

  return (
    <S.HeaderContainer>
      <S.Row>
        <S.HeaderText>Test Results</S.HeaderText>
        <div>
          {Boolean(totalPassedSpecs) && (
            <S.HeaderDetail>
              <S.HeaderDot $passed />
              {`${totalPassedSpecs} ${singularOrPlural('spec', totalPassedSpecs)} passed`}
            </S.HeaderDetail>
          )}

          {Boolean(totalFailedSpecs) && (
            <S.HeaderDetail>
              <S.HeaderDot $passed={false} />
              {`${totalFailedSpecs} ${singularOrPlural('spec', totalFailedSpecs)} failed`}
            </S.HeaderDetail>
          )}
        </div>
      </S.Row>

      <S.PrimaryButton
        data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.AddTestSpec)}
        data-cy="add-test-spec-button"
        onClick={handleAddTestSpecOnClick}
      >
        Add Test Spec
      </S.PrimaryButton>
    </S.HeaderContainer>
  );
};

export default Header;
