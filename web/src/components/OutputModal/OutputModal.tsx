import {Form} from 'antd';
import {useEffect} from 'react';
import {TOutput} from 'types/Output.types';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import SpanSelectors from 'selectors/Span.selectors';
import TestRunSelectors from 'selectors/TestRun.selectors';
import useValidateOutput from './hooks/useValidateOutput';
import * as S from './OutputModal.styled';
import OutputModalFooter from './OutputModalFooter';
import OutputModalForm from './OutputModalForm';

interface IProps {
  isOpen: boolean;
  onClose(): void;
  draftOutput?: TOutput;
  onSubmit(values: TOutput): void;
}

const OutputModal = ({isOpen, onClose, draftOutput, onSubmit}: IProps) => {
  const isEditing = Boolean(draftOutput);
  const [form] = Form.useForm<TOutput>();

  const {
    run: {id: runId},
  } = useTestRun();
  const {
    test: {id: testId},
  } = useTest();
  const isTraceSource = Form.useWatch('source', form) === 'trace';

  const spanIdList = useAppSelector(SpanSelectors.selectMatchedSpans);
  const {isValid, onValidate} = useValidateOutput({spanIdList});

  const traceAttributeList = useAppSelector(state =>
    AssertionSelectors.selectAttributeList(state, testId, runId, spanIdList)
  );
  const triggerAttributeList = useAppSelector(state =>
    TestRunSelectors.selectResponseAttributeList(state, testId, runId)
  );

  const attributeList = isTraceSource ? traceAttributeList : triggerAttributeList;

  useEffect(() => {
    if (draftOutput && isOpen) form.setFieldsValue(draftOutput);
    if (!isOpen || !draftOutput) form.resetFields();
  }, [draftOutput, form, isOpen]);

  useEffect(() => {
    onValidate(null, form.getFieldsValue());
    form.validateFields();
  }, [form, onValidate]);

  return isOpen ? (
    <S.Modal
      visible={isOpen}
      onCancel={onClose}
      footer={
        <OutputModalFooter
          isValid={isValid}
          isLoading={false}
          isEditing={isEditing}
          onCancel={onClose}
          onSave={() => form.submit()}
        />
      }
      title={<S.Title>{isEditing ? 'Edit Output' : 'Add Output'}</S.Title>}
      width={520}
    >
      <Form<TOutput>
        form={form}
        autoComplete="off"
        layout="vertical"
        onFinish={onSubmit}
        onValuesChange={onValidate}
        initialValues={{
          source: 'trigger',
          ...draftOutput,
        }}
      >
        <OutputModalForm
          form={form}
          spanIdList={spanIdList}
          attributeList={attributeList}
          testId={testId}
          runId={runId}
        />
      </Form>
    </S.Modal>
  ) : null;
};

export default OutputModal;
