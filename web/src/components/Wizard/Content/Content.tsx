import {Tabs} from 'antd';
import {IWizardStep} from 'types/Wizard.types';
import * as S from './Content.styled';
import Tab from './Tab';

interface IProps {
  activeStepId: string;
  isLoading: boolean;
  onChange(key: string): void;
  onNext(): void;
  steps: IWizardStep[];
}

const Content = ({activeStepId, isLoading, onChange, onNext, steps}: IProps) => (
  <S.Container>
    <Tabs
      type="card"
      activeKey={activeStepId}
      tabPosition="left"
      onTabClick={(key, event) => {
        event.preventDefault();
        onChange(key);
      }}
    >
      {steps.map((step, index) => {
        const {component: Component, id} = step;
        return (
          <Tabs.TabPane
            disabled={!step.isEnabled}
            tab={<Tab index={index + 1} isActive={id === activeStepId} step={step} />}
            key={id}
          >
            <Component isLoading={isLoading} onNext={onNext} />
          </Tabs.TabPane>
        );
      })}
    </Tabs>
  </S.Container>
);

export default Content;
