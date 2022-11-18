import {Form, FormInstance, Input, Tag} from 'antd';
import {noop} from 'lodash';

import Editor from 'components/Editor';
import useQuerySelector from 'components/TestSpecForm/hooks/useQuerySelector';
import {SupportedEditors} from 'constants/Editor.constants';
import {TTestOutput} from 'types/TestOutput.types';
import {singularOrPlural} from 'utils/Common';
import * as S from './OutputModal.styled';

interface IProps {
  form: FormInstance<TTestOutput>;
  runId: string;
  spanIdList: string[];
  testId: string;
}

const OutputModalForm = ({form, runId, spanIdList, testId}: IProps) => {
  useQuerySelector({
    form,
    runId,
    testId,
    onValidSelector: noop,
  });

  const selector = Form.useWatch('selector', form) || '';

  return (
    <S.InputContainer>
      <Form.Item
        data-cy="output-form-name"
        label="Name"
        name="name"
        rules={[{required: true, message: 'Please enter a name'}]}
        style={{marginBottom: 0}}
      >
        <Input />
      </Form.Item>

      <Form.Item
        data-cy="output-form-selector"
        label={
          <S.SelectorTitleContainer>
            <S.SelectorLabel>Selector</S.SelectorLabel>
            <Tag color="blue">{`${spanIdList.length} ${singularOrPlural('span', spanIdList.length)} selected`}</Tag>
          </S.SelectorTitleContainer>
        }
        name="selector"
        rules={[
          {required: true, message: 'Please enter a valid selector'},
          {
            message: 'Please select a single span',
            validator: async () => {
              if (spanIdList.length !== 1) throw new Error('Please select a single span');
            },
          },
        ]}
        style={{marginBottom: 0}}
      >
        <Editor basicSetup={{lineNumbers: true}} type={SupportedEditors.Selector} placeholder="" />
      </Form.Item>

      <Form.Item
        data-cy="output-form-value"
        name="value"
        rules={[{required: true, message: 'Please enter an attribute'}]}
        style={{margin: 0}}
      >
        <Editor
          type={SupportedEditors.Expression}
          placeholder="Attribute"
          context={{
            runId,
            testId,
            spanId: spanIdList[0],
            selector,
          }}
        />
      </Form.Item>
    </S.InputContainer>
  );
};

export default OutputModalForm;
