import {FormInstance, Select} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {uniqBy} from 'lodash';
import {useMemo} from 'react';
import AttributesTags from '../../constants/AttributesTags.json';
import {TAssertion} from '../../types/Assertion.types';
import {TSpanFlatAttribute} from '../../types/Span.types';
import {IValues} from './AssertionForm';
import {AssertionFormCheck} from './AssertionFormCheck';
import {
  OtelReference,
  useGetOTELSemanticConvertionAttributesInfo,
} from './hooks/useGetOTELSemanticConvertionAttributesInfo';

interface IProps {
  form: FormInstance<IValues>;
  fields: FormListFieldData[];
  assertionList: TAssertion[];

  add(): void;

  remove(name: number): void;

  attributeList: TSpanFlatAttribute[];
}

const AssertionFormCheckList: React.FC<IProps> = ({form, fields, add, remove, attributeList, assertionList}) => {
  const reference = useGetOTELSemanticConvertionAttributesInfo();
  const attributeOptionList = useMemo(() => {
    return uniqBy(attributeList, 'key').map(({key}) => (
      <Select.Option key={key} value={key}>
        {`${key} ${
          reference[key] || (AttributesTags as OtelReference)[key]
            ? reference[key]
              ? ` - ${reference[key].description}`
              : ` - ${(AttributesTags as OtelReference)[key].description}`
            : ''
        }`}
      </Select.Option>
    ));
  }, [attributeList, reference]);

  return (
    <>
      {fields.map(({key, name, ...field}, index) => {
        return (
          <AssertionFormCheck
            key={key}
            reference={reference}
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
    </>
  );
};

export default AssertionFormCheckList;
