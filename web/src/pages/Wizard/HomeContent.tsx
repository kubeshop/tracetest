import {Collapse} from 'antd';
import ContentResources from 'components/Resources/Content';
import HeaderResources from 'components/Resources/Header';
import Content from 'components/Wizard/Content';
import Header from 'components/Wizard/Header';
import {withCustomization} from 'providers/Customization';
import {useWizard} from 'providers/Wizard';
import * as S from './Wizard.styled';

const HomeContent = () => {
  const {activeStepId, isLoading, onGoTo, onNext, steps} = useWizard();
  const completedSteps = steps.filter(({state}) => state === 'completed').length;
  const isWizardComplete = !!steps.length && completedSteps === steps.length;

  return (
    <S.Container>
      <S.Header>
        <S.Title>Welcome to Tracetest!</S.Title>
        <S.Text>
          Here&apos;s a guide to get started and help you test your modern applications with OpenTelemetry.
        </S.Text>
      </S.Header>

      {!!steps.length && (
        <>
          <Collapse defaultActiveKey={!isWizardComplete ? 'wizard' : ''}>
            <Collapse.Panel
              header={<Header activeStep={completedSteps} totalCompleteSteps={steps.length} />}
              key="wizard"
            >
              <Content
                activeStepId={activeStepId}
                isLoading={isLoading}
                onChange={onGoTo}
                onNext={onNext}
                steps={steps}
              />
            </Collapse.Panel>
          </Collapse>

          <Collapse defaultActiveKey={isWizardComplete ? 'resources' : ''}>
            <Collapse.Panel header={<HeaderResources />} key="resources">
              <ContentResources />
            </Collapse.Panel>
          </Collapse>
        </>
      )}
    </S.Container>
  );
};

export default withCustomization(HomeContent, 'homeContent');
