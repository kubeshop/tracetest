import {useTour} from '@reactour/tour';
import {Button, Form, FormInstance, Modal, Typography} from 'antd';
import {useCallback, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useTheme} from 'styled-components';

import {IBasicDetailsValues} from 'components/CreateTestPlugins/Default/steps/BasicDetails/BasicDetails';
import {IRequestDetailsValues} from 'components/CreateTestPlugins/Rest/steps/RequestDetails/RequestDetails';
import {useEditTestMutation, useRunTestMutation} from 'redux/apis/TraceTest.api';
import {TRawHTTPRequest, TTest} from 'types/Test.types';
import TestDefinitionService from 'services/TestDefinition.service';
import {TriggerTypes} from 'constants/Test.constants';
import EditTestForm, {FORM_ID} from './EditTestForm';

export type TEditTest = IRequestDetailsValues & IBasicDetailsValues;
export type TEditTestForm = FormInstance<TEditTest>;

interface IProps {
  isOpen: boolean;
  onClose(): void;
  test: TTest;
}

const EditTestModal = ({onClose, isOpen, test}: IProps) => {
  const theme = useTheme();
  const [isFormValid, setIsFormValid] = useState(true);
  const navigate = useNavigate();
  const {setIsOpen} = useTour();
  const [editTest, {isLoading: isLoadingCreateTest}] = useEditTestMutation();
  const [runTest, {isLoading: isLoadingRunTest}] = useRunTestMutation();
  const [form] = Form.useForm<TEditTest>();

  const handleOnSubmit = useCallback(
    async (values: TEditTest) => {
      const request: TRawHTTPRequest = {
        url: values.url,
        method: values.method,
        body: values.body,
        headers: values.headers,
      };

      await editTest({
        test: {
          name: values.name,
          description: values.description,
          serviceUnderTest: {
            triggerType: TriggerTypes.http,
            triggerSettings: {
              http: request,
            },
          },
          definition: {
            definitions: test.definition.definitionList.map(def => TestDefinitionService.toRaw(def)),
          },
        },
        testId: test.id,
      }).unwrap();

      const run = await runTest({testId: test.id}).unwrap();
      onClose();
      setIsOpen(false);

      navigate(`/test/${test.id}/run/${run.id}`);
    },
    [editTest, navigate, onClose, runTest, setIsOpen, test.definition, test.id]
  );

  const handleOnCancel = useCallback(() => {
    setIsFormValid(false);
    form.resetFields();
    onClose();
  }, [form, onClose]);

  const footer = (
    <>
      <Button ghost onClick={handleOnCancel} type="primary">
        Cancel
      </Button>
      <Button
        data-cy="edit-test-submit"
        disabled={!isFormValid}
        form={FORM_ID}
        loading={isLoadingCreateTest || isLoadingRunTest}
        onClick={() => form.submit()}
        type="primary"
      >
        Save
      </Button>
    </>
  );

  return (
    <Modal
      bodyStyle={{backgroundColor: theme.color.background, overflowY: 'auto'}}
      footer={footer}
      width="40%"
      onCancel={handleOnCancel}
      title={
        <Typography.Title level={2} style={{marginBottom: 0}}>
          Edit Test - {test.name}
        </Typography.Title>
      }
      visible={isOpen}
    >
      <EditTestForm form={form} test={test} onSubmit={handleOnSubmit} onValidation={setIsFormValid} />
    </Modal>
  );
};

export default EditTestModal;
