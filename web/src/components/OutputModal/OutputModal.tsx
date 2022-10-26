import {Form} from 'antd';
import {useEffect} from 'react';

import {useAppSelector} from 'redux/hooks';
import {selectTestOutputByIndex} from 'redux/testOutputs/selectors';
import SpanSelectors from 'selectors/Span.selectors';
import {TTestOutput} from 'types/TestOutput.types';
import useValidateOutput from './hooks/useValidateOutput';
import * as S from './OutputModal.styled';
import OutputModalFooter from './OutputModalFooter';
import OutputModalForm from './OutputModalForm';

interface IProps {
  index: number;
  isOpen: boolean;
  onClose(): void;
  onSubmit(values: TTestOutput, isEditing: boolean): void;
  runId: string;
  testId: string;
}

const OutputModal = ({index, isOpen, onClose, onSubmit, runId, testId}: IProps) => {
  const [form] = Form.useForm<TTestOutput>();
  const output = useAppSelector(state => selectTestOutputByIndex(state, index));
  const spanIdList = useAppSelector(SpanSelectors.selectMatchedSpans);
  const {isValid, onValidate} = useValidateOutput({spanIdList});
  const isEditing = index !== -1;

  useEffect(() => {
    if (output && isOpen) form.setFieldsValue(output);
    if (!isOpen || !output) form.resetFields();
  }, [output, form, isOpen]);

  useEffect(() => {
    if (isOpen && Boolean(form.getFieldValue('selector'))) {
      onValidate(null, form.getFieldsValue());
      form.validateFields();
    }
  }, [form, isOpen, onValidate]);

  return (
    <S.Modal
      footer={
        <OutputModalFooter isValid={isValid} isEditing={isEditing} onCancel={onClose} onSave={() => form.submit()} />
      }
      onCancel={onClose}
      title={<S.Title>{isEditing ? 'Edit Output' : 'Add Output'}</S.Title>}
      visible={isOpen}
      width={520}
    >
      <Form<TTestOutput>
        autoComplete="off"
        form={form}
        layout="vertical"
        onFinish={values => onSubmit(values, isEditing)}
        onValuesChange={onValidate}
      >
        <OutputModalForm form={form} runId={runId} spanIdList={spanIdList} testId={testId} />
      </Form>
    </S.Modal>
  );
};

export default OutputModal;
