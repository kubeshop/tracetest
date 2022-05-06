import {useCallback} from 'react';
import {useTour} from '@reactour/tour';
import {Modal, Form, Button} from 'antd';
import {useCreateTestMutation, useRunTestMutation} from 'redux/apis/Test.api';
import './CreateTestModal.styled.ts';
import {useNavigate} from 'react-router-dom';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {Steps} from '../GuidedTour/homeStepList';
import CreateTestAnalyticsService from '../../services/Analytics/CreateTestAnalytics.service';
import * as S from './CreateTestModal.styled';
import CreateTestForm, {ICreateTestValues} from './CreateTestForm';

interface IProps {
  visible: boolean;
  onClose: () => void;
}

const {onCreateTestFormSubmit} = CreateTestAnalyticsService;

const CreateTestModal = ({visible, onClose}: IProps): JSX.Element => {
  const navigate = useNavigate();
  const {setIsOpen} = useTour();
  const [createTest, {isLoading: isLoadingCreateTest}] = useCreateTestMutation();
  const [runTest, {isLoading: isLoadingRunTest}] = useRunTestMutation();

  const [form] = Form.useForm<ICreateTestValues>();

  const onCreate = useCallback(
    async (values: ICreateTestValues) => {
      const headers = values.headersList
        .filter((i: {checked: boolean}) => i.checked)
        .map(({key, value}: {key: string; value: string}) => ({key, value}));
      const newTest = await createTest({
        name: values.name,
        serviceUnderTest: {
          request: {url: values.url, method: values.method, body: values.body, headers},
        },
      }).unwrap();
      const newTestRunResult = await runTest(newTest.testId).unwrap();
      onClose();
      setIsOpen(false);
      navigate(`/test/${newTest.testId}?resultId=${newTestRunResult.resultId}`);
    },
    [createTest, navigate, onClose, runTest, setIsOpen]
  );

  const renderActionButtons = () => {
    return (
      <>
        <Button type="ghost" htmlType="button" onClick={onClose}>
          Cancel
        </Button>

        <Button
          type="primary"
          form="newTest"
          onClick={() => {
            onCreateTestFormSubmit();
            form.submit();
          }}
          loading={isLoadingCreateTest || isLoadingRunTest}
          data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Run)}
          data-cy="create-test-submit"
        >
          Run Test
        </Button>
      </>
    );
  };

  return (
    <S.Wrapper>
      <Modal
        closable
        title="Create New Test"
        visible={visible}
        footer={renderActionButtons()}
        onCancel={onClose}
        wrapClassName="test-modal"
      >
        <div style={{display: 'flex'}} data-cy="create-test-modal">
          <CreateTestForm onSubmit={onCreate} form={form} />
        </div>
      </Modal>
    </S.Wrapper>
  );
};

export default CreateTestModal;
