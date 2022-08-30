import {FormInstance, Input, Select} from 'antd';
import {TAssertion} from '../../types/Assertion.types';
import {enumKeys} from '../../utils/Common';
import {IValues} from './TestSpecForm';
import * as S from './TestSpecForm.styled';

export enum Duration {
  nanoseconds = 'ns',
  microseconds = 'μs',
  milliseconds = 'ms',
  seconds = 's',
  minutes = 'm',
  hours = 'h',
}

interface DurationState {
  value?: string;
  duration?: Duration;
}

function setDurationFieldValue(form: FormInstance<IValues>, index: number, state: DurationState) {
  form.setFieldsValue({
    assertionList: form
      .getFieldsValue()
      .assertionList?.map((as, i) => (i === index ? {...as, expected: `${state.value}${state.duration}`} : as)),
  });
}

interface IProps {
  assertion: TAssertion;
  form: FormInstance<IValues>;
  index: number;
}

const durations = enumKeys(Duration).map(key => ({label: key, value: Duration[key]}));

export const DurationFields = ({form, index, assertion}: IProps) => {
  const state: DurationState = {
    value: assertion?.expected?.match(/(\d+)/)?.[0],
    duration: (assertion?.expected?.match(/[a-zA-Z]+/g)?.[0] || 'ms') as Duration,
  };
  return (
    <div style={{display: 'flex'}}>
      <Input
        value={state.value}
        onChange={e => setDurationFieldValue(form, index, {duration: state.duration, value: e.target.value})}
        placeholder="Expected Value"
        type="number"
        data-cy="duration-value"
      />
      <S.Select
        onChange={e => setDurationFieldValue(form, index, {value: state.value, duration: e as Duration})}
        value={state.duration}
        style={{margin: 0}}
        placeholder="Assertion Type"
        data-cy="duration"
      >
        {durations.map(({label, value}) => (
          <Select.Option data-cy={`duration-unit-${value}`} key={label} value={value}>
            {label}
          </Select.Option>
        ))}
      </S.Select>
    </div>
  );
};
