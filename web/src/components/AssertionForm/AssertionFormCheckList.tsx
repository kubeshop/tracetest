import {FormInstance, Select} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {uniqBy} from 'lodash';
import {useMemo} from 'react';
import {TAssertion} from '../../types/Assertion.types';
import {TSpanFlatAttribute} from '../../types/Span.types';
import {IValues} from './AssertionForm';
import {AssertionFormCheck} from './AssertionFormCheck';

interface IProps {
  form: FormInstance<IValues>;
  fields: FormListFieldData[];
  assertionList: TAssertion[];

  add(): void;

  remove(name: number): void;

  attributeList: TSpanFlatAttribute[];
}

const AssertionFormCheckList: React.FC<IProps> = ({form, fields, add, remove, attributeList, assertionList}) => {
  const attributeOptionList = useMemo(
    () =>
      uniqBy(attributeList, 'key').map(({key}) => (
        <Select.Option key={key} value={key}>
          {key}
        </Select.Option>
      )),
    [attributeList]
  );

  return (
    <div style={{display: 'flex', flexDirection: 'row'}}>
      {fields.map(({key, name, ...field}, index) => {
        return (
          <AssertionFormCheck
            key={key}
            form={form}
            add={add}
            remove={remove}
            field={field}
            attributeOptionList={attributeOptionList}
            name={name}
            index={index}
            length={fields.length}
            assertionList={assertionList}
          />
        );
      })}
    </div>
  );
};

export default AssertionFormCheckList;
