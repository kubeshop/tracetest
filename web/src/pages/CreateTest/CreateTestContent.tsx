import {useEffect} from 'react';
import {useNavigate} from 'react-router-dom';
import CreateTestHeader from 'components/CreateTestHeader';
import CreateTestSteps from 'components/CreateTestSteps';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import * as S from './CreateTest.styled';

const CreateTestContent: React.FC = () => {
  const navigate = useNavigate();
  const {stepList, activeStep, pluginName, onGoTo, onReset, isLoading} = useCreateTest();

  useEffect(() => {
    return () => {
      onReset();
    };
  }, []);

  return (
    <S.Wrapper>
      <CreateTestHeader onBack={() => navigate('/')} />
      <CreateTestSteps
        isLoading={isLoading}
        onGoTo={onGoTo}
        stepList={stepList}
        selectedKey={activeStep}
        pluginName={pluginName}
      />
    </S.Wrapper>
  );
};

export default CreateTestContent;
