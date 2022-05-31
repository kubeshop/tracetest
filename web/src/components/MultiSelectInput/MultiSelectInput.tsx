import {useState, useCallback, useMemo} from 'react';
import {Select} from 'antd';
import {LabeledValue} from 'antd/lib/select';
import MultiSelectInputTag from './MultiSelectInputTag';

interface IProps {
  entryList: {
    name: string;
    items: {key: string; value: string}[];
  }[];
  onEntry(entry: string[]): void;
  onStepEntry?(name: string, entry: string): void;
  onDeselect(entry: string[]): void;
  onClear(): void;
  placeholder?: string;
  defaultValueList?: LabeledValue[];
}

export const SEPARATOR = '~~~';

const MultiSelectInput: React.FC<IProps> = ({
  entryList,
  onEntry,
  onStepEntry,
  onDeselect,
  onClear,
  placeholder,
  defaultValueList = [],
}) => {
  const [selectedItemList, setSelectedItemList] = useState<LabeledValue[]>(defaultValueList);
  const [step, setStep] = useState(0);
  const [entry, setEntry] = useState<string[]>([]);
  const [entryCount, setEntryCount] = useState(defaultValueList.length);
  const {items, name} = useMemo(() => entryList[step], [entryList, step]);

  const optionList = useMemo(
    () => (
      <Select.OptGroup label={name}>
        {items.map(({key: itemKey, value}) => (
          <Select.Option
            value={`${value}${SEPARATOR}${step}${SEPARATOR}${entryCount}`}
            key={`${entryCount}-${itemKey}`}
          >
            {itemKey}
          </Select.Option>
        ))}
      </Select.OptGroup>
    ),
    [entryCount, items, name, step]
  );

  const handleChange = useCallback(
    (list: LabeledValue[]) => {
      if (!list.length) {
        onClear();

        setEntry([]);
        setStep(0);
        return setSelectedItemList(list);
      }

      setSelectedItemList(list);
      const {value: currentValue} = list[list.length - 1];
      const [value] = String(currentValue).split(SEPARATOR);
      const isLastStep = step === entryList.length - 1;

      if (isLastStep) {
        setStep(0);
        setEntry([]);
        setEntryCount(entryCount + 1);
        onEntry([...entry, value]);
      } else {
        if (onStepEntry) onStepEntry(name.toLowerCase(), value);
        setEntry(prevEntry => [...prevEntry, value]);
        setStep(step + 1);
      }
    },
    [entry, entryCount, entryList.length, name, onClear, onEntry, onStepEntry, step]
  );

  const handleDeselect = useCallback(
    (entryNumber: number) => {
      setSelectedItemList(
        selectedItemList.filter(({value}) => {
          const [, , number] = String(value).split(SEPARATOR);

          return Number(number) !== entryNumber;
        })
      );

      const deselectedEntry = selectedItemList.reduce<string[]>((list, {value, label}) => {
        const [, , number] = String(value).split(SEPARATOR);

        return Number(number) === entryNumber ? list.concat(label as string) : list;
      }, []);
      onDeselect(deselectedEntry);
    },
    [onDeselect, selectedItemList]
  );

  return (
    <Select
      data-cy="multi-select-input"
      mode="multiple"
      allowClear
      onChange={handleChange}
      labelInValue
      placeholder={placeholder}
      value={selectedItemList}
      tagRender={props => (
        <MultiSelectInputTag {...props} onDeselect={handleDeselect} entryListCount={entryList.length} />
      )}
    >
      {optionList}
    </Select>
  );
};

export default MultiSelectInput;
