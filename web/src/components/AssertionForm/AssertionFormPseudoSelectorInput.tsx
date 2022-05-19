import {Input, Select} from 'antd';
import {capitalize, noop} from 'lodash';
import {useCallback} from 'react';
import {PseudoSelector} from '../../constants/Operator.constants';
import * as S from './AssertionForm.styled';

const optionList = Object.entries(PseudoSelector).map(([key, value]) => (
  <Select.Option key={key} value={value}>
    {capitalize(key)}
  </Select.Option>
));

interface IPseudoSelector {
  selector: PseudoSelector;
  number?: number;
}

interface IProps {
  value?: IPseudoSelector;
  onChange?(pseudoSelector: IPseudoSelector): void;
}

const AssertionFormPseudoSelectorInput: React.FC<IProps> = ({
  value: pseudoSelector = {
    selector: PseudoSelector.FIRST,
    number: 0,
  },
  onChange = noop,
}) => {
  const handleSelectorChange = useCallback(
    (selector: PseudoSelector) => {
      onChange({
        selector,
        number: pseudoSelector.number,
      });
    },
    [onChange, pseudoSelector.number]
  );

  const handleNumberChange = useCallback(
    event => {
      onChange({
        selector: pseudoSelector.selector,
        number: Number(event.target.value),
      });
    },
    [onChange, pseudoSelector.selector]
  );

  return (
    <S.PseudoSelector>
      <Select onChange={handleSelectorChange} value={pseudoSelector?.selector} style={{width: '79px'}}>
        {optionList}
      </Select>
      {pseudoSelector?.selector === PseudoSelector.NTH && (
        <Input placeholder="number" style={{width: '100px'}} type="number" onChange={handleNumberChange} />
      )}
    </S.PseudoSelector>
  );
};

export default AssertionFormPseudoSelectorInput;
