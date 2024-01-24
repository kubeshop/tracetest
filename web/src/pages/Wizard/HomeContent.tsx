import Header from 'components/Wizard/Header';
import Content from 'components/Wizard/Content';
import {withCustomization} from 'providers/Customization';
import {useWizard} from 'providers/Wizard';
import * as S from './Wizard.styled';

const HomeContent = () => {
  const {activeStepId, isLoading, onGoTo, onNext, steps} = useWizard();
  const completedSteps = steps.filter(({state}) => state === 'completed').length;

  return (
    <S.Container>
      <S.Header>
        <S.Title>Welcome to Tracetest!</S.Title>
        <S.Text>
          Here&apos;s a guide to get started and help you test your modern applications with OpenTelemetry.
        </S.Text>
      </S.Header>

      <S.Body>
        <Header activeStep={completedSteps} totalCompleteSteps={steps.length} />
        <Content activeStepId={activeStepId} isLoading={isLoading} onChange={onGoTo} onNext={onNext} steps={steps} />
      </S.Body>
    </S.Container>
  );
};

export default withCustomization(HomeContent, 'homeContent');
