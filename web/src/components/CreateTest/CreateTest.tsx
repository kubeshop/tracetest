import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {Empty, Form, Typography} from 'antd';
import {TDraftTest} from 'types/Test.types';
import {TriggerTypes} from 'constants/Test.constants';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {useState} from 'react';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import EditRequestDetails from 'components/EditTestForm/EditRequestDetails/EditRequestDetails';
import * as S from './CreateTest.styled';
import Header from './Header';

export const FORM_ID = 'create-test';

interface IProps {
  triggerType: TriggerTypes;
}

const CreateTest = ({triggerType}: IProps) => {
  const {onCreateTest} = useCreateTest();
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
      </Form>
    </S.Container>
  );
};

export default CreateTest;
