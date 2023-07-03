import {Button, Form, Input, Tag} from 'antd';
import {useEffect} from 'react';

import Editor from 'components/Editor';
import {EXPRESSIONS_DOCUMENTATION_URL, SELECTOR_LANGUAGE_CHEAT_SHEET_URL} from 'constants/Common.constants';
import {SupportedEditors} from 'constants/Editor.constants';
import {useAppSelector} from 'redux/hooks';
import SpanSelectors from 'selectors/Span.selectors';
import {singularOrPlural} from 'utils/Common';
import TestOutput from 'models/TestOutput.model';
import SelectorInput from './SelectorInput';
import * as S from './TestOutputForm.styled';

interface IProps {
  isEditing?: boolean;
  isLoading?: boolean;
  onCancel(): void;
  onSubmit(values: TestOutput, spanId?: string): void;
  onValidate(_: any, output: TestOutput): void;
  isValid: boolean;
  output?: TestOutput;
  runId: string;
  testId: string;
}

const TestOutputForm = ({
  isEditing = false,
  isLoading = false,
  onCancel,
  onSubmit,
  output,
  runId,
  testId,
  onValidate,
  isValid,
}: IProps) => {
  const [form] = Form.useForm<TestOutput>();
  const selector = Form.useWatch('selector', form) || '';
  const spanIdList = useAppSelector(SpanSelectors.selectMatchedSpans);

  useEffect(() => {
    if (form.getFieldValue('selector') && form.getFieldValue('value') && form.getFieldValue('name')) {
      onValidate(null, form.getFieldsValue());
      form.validateFields();
    }
  }, [form, onValidate]);

  useEffect(() => {
    if (!isEditing) {
      form.setFieldsValue({selector: output?.selector, value: output?.value, name: output?.name});
    }
  }, [form, isEditing, output?.name, output?.selector, output?.value]);

  useEffect(() => {
    return () => {
      onCancel();
    };
  }, []);

  return (
    <S.Container>
      <S.Title>{isEditing ? 'Edit Test Output' : 'Add Test Output'}</S.Title>

      <Form<TestOutput>
        autoComplete="off"
        form={form}
        initialValues={output}
        layout="vertical"
        name="testOutput"
        onFinish={values => onSubmit(values, spanIdList[0])}
        onValuesChange={onValidate}
      >
        <Form.Item hidden name="spanId" />
        <S.FormSection>
          <S.FormSectionHeaderSelector>
            <S.FormSectionRow1>
              <S.FormSectionTitle $noMargin>1. Select one Span</S.FormSectionTitle>
              <Tag color="blue">{`${spanIdList.length} ${singularOrPlural('span', spanIdList.length)} selected`}</Tag>
            </S.FormSectionRow1>
            <a href={SELECTOR_LANGUAGE_CHEAT_SHEET_URL} target="_blank">
              <S.ReadIcon /> SL cheat sheet
            </a>
          </S.FormSectionHeaderSelector>
          <S.FormSectionRow>
            <S.FormSectionText>Select only one span to extract the value from its attributes</S.FormSectionText>
          </S.FormSectionRow>
          <SelectorInput form={form} runId={runId} spanIdList={spanIdList} testId={testId} />
        </S.FormSection>

        <S.FormSection>
          <S.FormSectionTitle>2. Select the attribute</S.FormSectionTitle>
          <S.FormSectionRow>
            <S.FormSectionText>Choose one attribute from the selected span or use an </S.FormSectionText>
            <a href={EXPRESSIONS_DOCUMENTATION_URL} target="_blank">
              expression
            </a>
          </S.FormSectionRow>
          <Form.Item name="value" rules={[{required: true, message: 'Please enter an attribute or expression'}]}>
            <Editor
              context={{
                runId,
                testId,
                spanId: spanIdList[0],
                selector,
              }}
              placeholder="Attribute"
              type={SupportedEditors.Expression}
            />
          </Form.Item>
        </S.FormSection>

        <S.FormSection>
          <S.FormSectionTitle>3. Give it a name</S.FormSectionTitle>
          <S.FormSectionRow>
            <S.FormSectionText>Give your output a unique name</S.FormSectionText>
          </S.FormSectionRow>
          <Form.Item name="name" rules={[{required: true, message: 'Please enter a name'}]}>
            <Input />
          </Form.Item>
        </S.FormSection>

        <S.Footer>
          <Button data-cy="output-modal-cancel-button" onClick={onCancel}>
            Cancel
          </Button>
          <Button data-cy="output-save-button" disabled={!isValid} htmlType="submit" loading={isLoading} type="primary">
            Save Test Output
          </Button>
        </S.Footer>
      </Form>
    </S.Container>
  );
};

export default TestOutputForm;
