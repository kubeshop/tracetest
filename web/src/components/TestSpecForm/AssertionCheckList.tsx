import {PlusOutlined} from '@ant-design/icons';
import {FormInstance} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import CreateAssertionModalAnalyticsService from 'services/Analytics/CreateAssertionModalAnalytics.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import {AssertionCheck} from './AssertionCheck';
import {useGetOTELSemanticConventionAttributesInfo} from './hooks/useGetOTELSemanticConventionAttributesInfo';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';

interface IProps {
  add(): void;
  attributeList: TSpanFlatAttribute[];
  fields: FormListFieldData[];
  form: FormInstance<IValues>;
  remove(name: number): void;
  runId: number;
  testId: string;
  spanIdList: string[];
}

const AssertionCheckList = ({form, fields, add, remove, attributeList, runId, testId, spanIdList}: IProps) => {
  const reference = useGetOTELSemanticConventionAttributesInfo();

  return (
    <S.AssertionsContainer>
      {fields.map(({key, name, ...field}) => {
        return (
          <AssertionCheck
            key={key}
            form={form}
            remove={remove}
            field={field}
            attributeList={attributeList}
            name={name}
            reference={reference}
            runId={runId}
            testId={testId}
            spanIdList={spanIdList}
          />
        );
      })}
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
