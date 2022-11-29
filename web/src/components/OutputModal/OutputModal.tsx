import {Form} from 'antd';
import {useEffect} from 'react';

import {useAppSelector} from 'redux/hooks';
import SpanSelectors from 'selectors/Span.selectors';
import {TTestOutput} from 'types/TestOutput.types';
import useValidateOutput from './hooks/useValidateOutput';
import * as S from './OutputModal.styled';
import OutputModalFooter from './OutputModalFooter';
import OutputModalForm from './OutputModalForm';

interface IProps {
  isOpen: boolean;
  onClose(): void;
  onSubmit(values: TTestOutput, isEditing: boolean): void;
  runId: string;
  testId: string;
  output?: TTestOutput;
  isEditing?: boolean;
}

const OutputModal = ({isOpen, onClose, onSubmit, runId, testId, output, isEditing = false}: IProps) => {
  const [form] = Form.useForm<TTestOutput>();
  const spanIdList = useAppSelector(SpanSelectors.selectMatchedSpans);
  const {isValid, onValidate} = useValidateOutput({spanIdList});

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
        name="testOutput"
        onFinish={values => onSubmit(values, isEditing)}
        onValuesChange={onValidate}
      >
        <OutputModalForm form={form} runId={runId} spanIdList={spanIdList} testId={testId} />
      </Form>
    </S.Modal>
  );
};

export default OutputModal;
