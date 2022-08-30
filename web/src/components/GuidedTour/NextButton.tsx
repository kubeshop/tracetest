import {BtnFnProps} from '@reactour/tour/dist/types';
import {PrimaryButton} from 'components/TestResults/TestResults.styled';

export const NextButton: React.FC<BtnFnProps> = ({currentStep, setCurrentStep, stepsLength, setIsOpen}) => {
  const isLast = currentStep === stepsLength - 1;
  return (
    <PrimaryButton
      onClick={() => {
        if (isLast) {
          setIsOpen(false);
          setCurrentStep(0);
          return;
        }
        setCurrentStep(prevState => Number(prevState) + 1);
      }}
    >
      {isLast ? 'Ok' : 'Next'}
    </PrimaryButton>
  );
};
