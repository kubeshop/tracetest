import {Form, FormInstance, Input, Select} from 'antd';
import {noop} from 'lodash';
import {useMemo} from 'react';
import {TOutput} from 'types/Output.types';
import OutputService from '../../services/Output.service';
import {TSpanFlatAttribute} from '../../types/Span.types';
import {singularOrPlural} from '../../utils/Common';
import AdvancedEditor from '../AdvancedEditor';
import {AttributeField} from '../TestSpecForm/Fields/AttributeField';
import {useGetOTELSemanticConventionAttributesInfo} from '../TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import useQuerySelector from '../TestSpecForm/hooks/useQuerySelector';
import * as S from './OutputModal.styled';
import OutputModalValue from './OutputModalValue';

interface IProps {
  form: FormInstance<TOutput>;
  attributeList: TSpanFlatAttribute[];
  spanIdList: string[];
  runId: string;
  testId: string;
}

const OutputModalForm = ({form, runId, testId, attributeList, spanIdList}: IProps) => {
  const otelReference = useGetOTELSemanticConventionAttributesInfo();
  const isTraceSource = Form.useWatch('source', form) === 'trace';
  const currentAttribute = Form.useWatch('attribute', form);
  const currentRegex = Form.useWatch('regex', form);

  useQuerySelector({
    form,
    runId,
    testId,
    onValidSelector: noop,
  });

  const value = useMemo(
    () => OutputService.getValueFromAttributeList(attributeList, currentAttribute, currentRegex),
    [attributeList, currentAttribute, currentRegex]
  );

  return (
    <S.InputContainer>
      <Form.Item
        data-cy="output-form-source"
        label="Source"
        name="source"
        rules={[{required: true, message: 'Please enter the output source'}]}
        style={{marginBottom: 0}}
      >
        <Select placeholder="From where the output is going to be extracted">
          <Select.Option value="trigger">Trigger</Select.Option>
          <Select.Option value="trace">Trace</Select.Option>
        </Select>
      </Form.Item>
      {isTraceSource && (
        <Form.Item
          data-cy="output-form-selector"
          label={
            <S.SelectorTitleContainer>
              <S.SelectorLabel>Selector</S.SelectorLabel>
              <S.SelectedCount>
                Selecting {spanIdList.length} {singularOrPlural('span', spanIdList.length)}
              </S.SelectedCount>
            </S.SelectorTitleContainer>
          }
          name="selector"
          style={{marginBottom: 0}}
          rules={[
            {required: true, message: 'Please enter a valid selector'},
            {
              message: 'Please select a single span',
              validator: async () => {
                if (spanIdList.length !== 1) throw new Error('Please select a single span');
              },
            },
          ]}
        >
          <AdvancedEditor lineNumbers runId={runId} testId={testId} />
        </Form.Item>
      )}
      <AttributeField
        data-cy="output-form-attribute"
        label="Attribute"
        name="attribute"
        style={{marginBottom: 0}}
        rules={[{required: true, message: 'Please enter an attribute'}]}
        attributeList={attributeList}
        reference={otelReference}
      />
      <Form.Item data-cy="output-form-regex" label="Regex" name="regex" style={{marginBottom: 0}}>
        <Input />
      </Form.Item>
      <S.ValueText>
        Value: <OutputModalValue value={value} />
      </S.ValueText>
    </S.InputContainer>
  );
};

export default OutputModalForm;
