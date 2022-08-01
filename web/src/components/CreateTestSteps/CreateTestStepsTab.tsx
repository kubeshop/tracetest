import {Typography} from 'antd';
import {TStepStatus} from '../../types/Plugins.types';
import * as S from './CreateTestSteps.styled';

interface IProps {
  text: string;
  status?: TStepStatus;
}

const CreateTestStepsTab = ({text, status}: IProps) => {
  return (
    <S.CreateTestStepsTab>
      <Typography.Text>{text}</Typography.Text>
      <S.StatusIcon $status={status} />
    </S.CreateTestStepsTab>
  );
};

export default CreateTestStepsTab;
