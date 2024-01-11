import {Tabs} from 'antd';
import {IWizardStep} from 'types/Wizard.types';
import * as S from './Steps.styled';
import StepTab from './StepTab';

interface IProps {
  activeStepId: string;
  componentFactory({step}: {step: IWizardStep}): React.ReactElement;
  steps: IWizardStep[];
}

const Steps = ({activeStepId, componentFactory: ComponentFactory, steps}: IProps) => (
  <S.Container>
    <Tabs
      type="card"
      activeKey={activeStepId}
      tabPosition="left"
      onTabClick={(key, event) => {
        event.preventDefault();
        // onGoTo(key);
      }}
    >
      {steps.map((step, index) => (
        <Tabs.TabPane tab={<StepTab index={index + 1} isActive={step.id === activeStepId} step={step} />} key={step.id}>
          <ComponentFactory step={step} />
        </Tabs.TabPane>
      ))}
    </Tabs>
  </S.Container>
);

export default Steps;
