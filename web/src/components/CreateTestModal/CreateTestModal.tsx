import {useTour} from '@reactour/tour';
import {Button, Form, Modal, Typography} from 'antd';

import {Steps} from 'components/GuidedTour/homeStepList';
import {useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {useCreateTestMutation, useRunTestMutation} from 'redux/apis/TraceTest.api';
import CreateTestAnalyticsService from 'services/Analytics/CreateTestAnalytics.service';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {Request} from '../../types/Common.types';
import CreateTestForm, {FORM_ID, ICreateTestValues} from './CreateTestForm';

interface IProps {
  visible: boolean;
  onClose: () => void;
}

const {onCreateTestFormSubmit} = CreateTestAnalyticsService;

const CreateTestModal = ({onClose, visible}: IProps) => {
  const [isFormValid, setIsFormValid] = useState(false);
  const [selectedDemo, setSelectedDemo] = useState('');
  const navigate = useNavigate();
  const {setIsOpen} = useTour();
  const [createTest, {isLoading: isLoadingCreateTest}] = useCreateTestMutation();
  const [runTest, {isLoading: isLoadingRunTest}] = useRunTestMutation();
  const [form] = Form.useForm<ICreateTestValues>();

  const handleOnSubmit = async (values: ICreateTestValues) => {
    let request: Request = {
      url: values.url,
      method: values.method,
      body: values.body,
      headers: values.headers,
    };
    if (values.auth?.type !== null) {
      request = {...request, auth: values.auth};
    }
    const test = await createTest({
      name: values.name,
      serviceUnderTest: {
        request,
      },
    }).unwrap();
    const run = await runTest({testId: test.id}).unwrap();
    onClose();
    setIsOpen(false);
    navigate(`/test/${test.id}/run/${run.id}`);
  };

  const handleOnCancel = () => {
    setIsFormValid(false);
    setSelectedDemo('');
    form.resetFields();
    onClose();
  };

  const footer = (
    <>
      <Button ghost onClick={handleOnCancel} type="primary">
        Cancel
      </Button>
      <Button
        data-cy="create-test-submit"
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Run)}
        disabled={!isFormValid}
        form={FORM_ID}
        loading={isLoadingCreateTest || isLoadingRunTest}
        onClick={() => {
          onCreateTestFormSubmit();
          form.submit();
        }}
        type="primary"
      >
        Create
      </Button>
    </>
  );

  return (
    <Modal
      bodyStyle={{backgroundColor: '#FBFBFF', maxHeight: 438, overflowY: 'auto'}}
      footer={footer}
      onCancel={handleOnCancel}
      title={
        <Typography.Title level={5} style={{marginBottom: 0}}>
          Create Test
        </Typography.Title>
      }
      visible={visible}
    >
      <CreateTestForm
        form={form}
        onSelectDemo={setSelectedDemo}
        onSubmit={handleOnSubmit}
        onValidation={setIsFormValid}
        selectedDemo={selectedDemo}
      />
    </Modal>
  );
};

export default CreateTestModal;
