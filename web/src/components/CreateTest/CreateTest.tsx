import {Empty, Typography} from 'antd';
import {TriggerTypes} from 'constants/Test.constants';
import useShortcut from 'components/TestPlugins/hooks/useShortcut';
import FormFactory from 'components/TestPlugins/FormFactory';
import CreateTestAnalytics from 'services/Analytics/CreateTest.service';
import * as S from './CreateTest.styled';
import Header from './Header';

interface IProps {
  isLoading: boolean;
  isValid: boolean;
  triggerType: TriggerTypes;
}

const CreateTest = ({isLoading, isValid, triggerType}: IProps) => {
  useShortcut();

  return (
    <S.Container>
      <Header
        isLoading={isLoading}
        isValid={isValid}
        onRunTest={() => CreateTestAnalytics.onFormSubmit()}
        triggerType={triggerType}
      />

      <S.Body>
        <S.SectionLeft>
          <FormFactory type={triggerType} />
        </S.SectionLeft>

        <S.SectionRight>
          <Typography.Title level={2}>Response</Typography.Title>
          <S.EmptyContainer>
            <Empty
              description="Enter the trigger details and click run to get a response"
              image={Empty.PRESENTED_IMAGE_SIMPLE}
            />
          </S.EmptyContainer>
        </S.SectionRight>
      </S.Body>
    </S.Container>
  );
};

export default CreateTest;
