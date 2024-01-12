import {Tabs} from 'antd';
import {IWizardStep} from 'types/Wizard.types';
import * as S from './Content.styled';
import Tab from './Tab';

interface IProps {
  activeStepId: string;
  steps: IWizardStep[];
  onChange(key: string): void;
}

const Content = ({activeStepId, steps, onChange}: IProps) => (
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
          <Tabs.TabPane tab={<Tab index={index + 1} isActive={id === activeStepId} step={step} />} key={id}>
            <Component />
          </Tabs.TabPane>
        );
      })}
    </Tabs>
  </S.Container>
);

export default Content;
