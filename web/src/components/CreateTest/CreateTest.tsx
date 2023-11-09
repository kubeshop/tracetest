import {Empty, Form, Typography} from 'antd';
import {TDraftTest} from 'types/Test.types';
import {TriggerTypes} from 'constants/Test.constants';
import EditRequestDetails from 'components/EditTestForm/EditRequestDetails/EditRequestDetails';
import * as S from './CreateTest.styled';
import Header from './Header';

export const FORM_ID = 'create-test';

interface IProps {
  triggerType: TriggerTypes;
  isValid: boolean;
}

const CreateTest = ({triggerType, isValid}: IProps) => {
  const form = Form.useFormInstance<TDraftTest>();

  return (
    <S.Container>
      <Header triggerType={triggerType} isValid={isValid} />

      <S.Body>
        <S.SectionLeft>
          <EditRequestDetails form={form} type={triggerType} />
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
