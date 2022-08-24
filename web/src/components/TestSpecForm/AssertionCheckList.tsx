import {PlusOutlined} from '@ant-design/icons';
import {FormInstance, Select} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {uniqBy} from 'lodash';
import {useMemo} from 'react';
import {TAssertion} from '../../types/Assertion.types';
import {TSpanFlatAttribute} from '../../types/Span.types';
import {IValues} from './TestSpecForm';
import {AssertionCheck} from './AssertionCheck';
import * as S from './TestSpecForm.styled';
import CreateAssertionModalAnalyticsService from '../../services/Analytics/CreateAssertionModalAnalytics.service';

interface IProps {
  form: FormInstance<IValues>;
  fields: FormListFieldData[];
  assertionList: TAssertion[];

  add(): void;

  remove(name: number): void;

  attributeList: TSpanFlatAttribute[];
}

const AssertionCheckList: React.FC<IProps> = ({form, fields, add, remove, attributeList, assertionList}) => {
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
    <S.AssertionsContainer>
      <S.CheckContainer>
        {fields.map(({key, name, ...field}, index) => {
          return (
            <AssertionCheck
              key={key}
              form={form}
              remove={remove}
              field={field}
              attributeOptionList={attributeOptionList}
              name={name}
              index={index}
              assertionList={assertionList}
            />
          );
        })}
      </S.CheckContainer>
      <S.AddCheckButton
        icon={<PlusOutlined />}
        onClick={() => {
          CreateAssertionModalAnalyticsService.onAddCheck();
          add();
        }}
        data-cy="add-assertion-form-add-check"
      >
        Add new
      </S.AddCheckButton>
    </S.AssertionsContainer>
  );
};

export default AssertionCheckList;
