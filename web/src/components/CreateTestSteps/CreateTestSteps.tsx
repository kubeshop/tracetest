import {Tabs} from 'antd';
import {SupportedPlugins} from 'constants/Plugins.constants';
import {ICreateTestStep} from 'types/Plugins.types';
import LoadingSpinner from 'components/LoadingSpinner';
import CreateTestStepFactory from './CreateTestStepFactory';
import * as S from './CreateTestSteps.styled';
import CreateTestStepsTab from './CreateTestStepsTab';

interface IProps {
  isLoading: boolean;
  stepList: ICreateTestStep[];
  selectedKey: string;
  pluginName: SupportedPlugins;
  onGoTo(stepId: string): void;
}

const CreateTestSteps = ({isLoading, stepList, selectedKey, pluginName, onGoTo}: IProps) => {
  return (
    <S.CreateTestSteps>
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
          <Tabs.TabPane tab={<CreateTestStepsTab text={step.name} status={step.status} />} key={step.id}>
            {isLoading ? (
              <S.LoadingSpinnerContainer>
                <LoadingSpinner />
              </S.LoadingSpinnerContainer>
            ) : (
              <CreateTestStepFactory pluginName={pluginName} step={step} />
            )}
          </Tabs.TabPane>
        ))}
      </Tabs>
    </S.CreateTestSteps>
  );
};

export default CreateTestSteps;
