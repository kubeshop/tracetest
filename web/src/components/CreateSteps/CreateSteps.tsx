import {Tabs} from 'antd';
import {ICreateTestStep} from 'types/Plugins.types';
import LoadingSpinner from 'components/LoadingSpinner';
import * as S from './CreateSteps.styled';
import CreateTestStepsTab from './CreateStepsTab';

interface IProps {
  isLoading: boolean;
  stepList: ICreateTestStep[];
  selectedKey: string;
  componentFactory({step}: {step: ICreateTestStep}): React.ReactElement;
  onGoTo(stepId: string): void;
  mode: string;
}

const CreateTestSteps = ({isLoading, stepList, selectedKey, componentFactory: ComponentFactory, onGoTo, mode}: IProps) => {
  return (
    <S.CreateTestSteps data-cy={`create-test-steps-${mode}`}>
      <S.ProgressLine $stepCount={stepList.length} />
      <Tabs
        type="card"
        activeKey={selectedKey}
        tabPosition="left"
        onTabClick={(key, event) => {
          event.preventDefault();
          onGoTo(key);
        }}
      >
        {stepList.map(step => (
          <Tabs.TabPane
            tab={<CreateTestStepsTab text={step.name} status={step.status} isActive={step.id === selectedKey} />}
            key={step.id}
          >
            {isLoading ? (
              <S.LoadingSpinnerContainer>
                <LoadingSpinner />
              </S.LoadingSpinnerContainer>
            ) : (
              <ComponentFactory step={step} />
            )}
          </Tabs.TabPane>
        ))}
      </Tabs>
    </S.CreateTestSteps>
  );
};

export default CreateTestSteps;
