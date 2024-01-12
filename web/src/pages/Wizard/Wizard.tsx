import withAnalytics from 'components/WithAnalytics/WithAnalytics';
import Header from 'components/Wizard/Header';
import StepFactory from 'components/Wizard/StepFactory';
import Steps from 'components/Wizard/Steps';
import {useWizard} from 'providers/Wizard';
import * as S from './Wizard.styled';

const Wizard = () => {
  const {activeStep, activeStepId, steps} = useWizard();

  return (
    <S.Container>
      <S.Header>
        <S.Title>Welcome to Tracetest!</S.Title>
        <S.Text>
          Here&apos;s a guide to get started and help you test your modern applications with OpenTelemetry.
        </S.Text>
      </S.Header>

      <S.Body>
        <Header activeStep={activeStep} totalSteps={steps.length} />
        <Steps activeStepId={activeStepId} componentFactory={StepFactory} steps={steps} />
      </S.Body>
    </S.Container>
  );
};

export default withAnalytics(Wizard, 'wizard');
