import {useNavigate} from 'react-router-dom';
import CreateTestHeader from 'components/CreateTestHeader';
import CreateTestSteps from 'components/CreateTestSteps';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import * as S from './CreateTest.styled';

const CreateTestContent: React.FC = () => {
  const navigate = useNavigate();
  const {stepList, activeStep, pluginName, onGoTo} = useCreateTest();

  return (
    <S.Wrapper>
      <CreateTestHeader onBack={() => navigate('/')} />
      <CreateTestSteps onGoTo={onGoTo} stepList={stepList} selectedKey={activeStep} pluginName={pluginName} />
    </S.Wrapper>
  );
};

export default CreateTestContent;
