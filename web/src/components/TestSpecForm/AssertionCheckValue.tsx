import {Dropdown, Form, FormInstance, Menu} from 'antd';
import {uniq} from 'lodash';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {useCallback, useMemo} from 'react';
import {SupportedEditors} from 'constants/Editor.constants';
import AssertionService from 'services/Assertion.service';
import {TStructuredAssertion} from 'types/Assertion.types';
import {TResolveExpressionContext} from 'types/Expression.types';
import Editor from '../Editor';
import {IValues} from './TestSpecForm';

interface IProps {
  form: FormInstance<IValues>;
  name: number;
  field: Pick<FormListFieldData, never>;
  valueList: string[];
  editorContext: TResolveExpressionContext;
}

const AssertionCheckValue = ({form, name, valueList, field, editorContext}: IProps) => {
  const onSelectedValue = useCallback(
    (value: string) => {
      const assertions = form.getFieldValue('assertions') as TStructuredAssertion[];

      form.setFieldsValue({
        assertions: assertions.map((assertion, i) =>
          i === name ? {...assertion, right: AssertionService.extractExpectedString(value) || ''} : assertion
        ),
      });
    },
    [form, name]
  );

  const menu = useMemo(
    () => (
      <Menu data-cy="assertion-check-value-menu">
        {uniq(valueList).map(value => (
          <Menu.Item key={value} onClick={() => onSelectedValue(value)}>
            {value}
          </Menu.Item>
        ))}
      </Menu>
    ),
    [onSelectedValue, valueList]
  );

  return (
    <Dropdown overlay={menu} trigger={['click']}>
      <Form.Item
        {...field}
        name={[name, 'right']}
        rules={[{required: true, message: 'Expected value is required'}]}
        data-cy="assertion-check-value"
        style={{margin: 0}}
      >
        <Editor type={SupportedEditors.Expression} placeholder="Expected Value" context={editorContext} />
      </Form.Item>
    </Dropdown>
  );
};

export default AssertionCheckValue;
