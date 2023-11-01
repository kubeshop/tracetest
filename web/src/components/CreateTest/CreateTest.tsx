import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {Empty, Form, Typography} from 'antd';
import {TDraftTest} from 'types/Test.types';
import {TriggerTypes} from 'constants/Test.constants';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {useState} from 'react';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import AllowButton, {Operation} from 'components/AllowButton';
import CreateButton from 'components/CreateButton';
import EditRequestDetails from 'components/EditTestForm/EditRequestDetails/EditRequestDetails';
import * as S from './CreateTest.styled';
import TriggerHeaderBar from './TriggerHeaderBar';

export const FORM_ID = 'create-test';

interface IProps {
  triggerType: TriggerTypes;
}

const CreateTest = ({triggerType}: IProps) => {
  const {onCreateTest, isLoading} = useCreateTest();
  const plugin = TriggerTypeToPlugin[triggerType];
  const [isValid, setIsValid] = useState(false);

  const onValidate = useValidateTestDraft({pluginName: plugin.name, setIsValid});
  const [form] = Form.useForm<TDraftTest>();

  const handleOnSubmit = async (values: TDraftTest) => {
    console.log(values);
    onCreateTest(values, plugin);
  };

  return (
    <S.Container>
      <Form<TDraftTest>
        autoComplete="off"
        data-cy="edit-test-modal"
        form={form}
        layout="vertical"
        name={FORM_ID}
        onFinish={handleOnSubmit}
        onValuesChange={onValidate}
      >
        <S.Header>
          <S.HeaderLeft>
            <TriggerHeaderBar form={form} type={triggerType} />
          </S.HeaderLeft>

          <S.HeaderRight>
            <AllowButton
              operation={Operation.Edit}
              block
              ButtonComponent={CreateButton}
              data-cy="edit-test-submit"
              disabled={!isValid}
              loading={isLoading}
              onClick={() => form.submit()}
              type="primary"
            >
              Run
            </AllowButton>
          </S.HeaderRight>
        </S.Header>

        <S.Body>
          <S.SectionLeft>
            <Typography.Title level={2}>Request Data</Typography.Title>

            <EditRequestDetails form={form} type={triggerType} />
          </S.SectionLeft>

          <S.SectionRight>
            <Typography.Title level={2}>Response Data</Typography.Title>
            <S.EmptyContainer>
              <Empty
                description="Enter the trigger details and click run to get a response"
                image={Empty.PRESENTED_IMAGE_SIMPLE}
              />
            </S.EmptyContainer>
          </S.SectionRight>
        </S.Body>
      </Form>
    </S.Container>
  );
};

export default CreateTest;
