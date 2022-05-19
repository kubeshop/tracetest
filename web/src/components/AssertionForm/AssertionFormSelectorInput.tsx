import {noop, uniqBy} from 'lodash';
import React, {useCallback, useMemo, useState} from 'react';
import {CompareOperator} from '../../constants/Operator.constants';
import {LOCATION_NAME} from '../../constants/Span.constants';
import {IItemSelector} from '../../types/Assertion.types';
import {ISpanFlatAttribute} from '../../types/Span.types';
import MultiSelectInput, {SEPARATOR} from '../MultiSelectInput/MultiSelectInput';
import * as S from './AssertionForm.styled';

type TItemSelectorDropdownProps = {
  attributeList: ISpanFlatAttribute[];
  value?: IItemSelector[];
  onChange?(selectorList: IItemSelector[]): void;
};

const operatorOptionList = [
  {
    key: '=',
    value: CompareOperator.EQUALS,
  },
  // this will be added when we add the query language parser
  // {
  //   key: 'contains',
  //   value: CompareOperator.CONTAINS,
  // },
];

const AssertionFormSelectorInput: React.FC<TItemSelectorDropdownProps> = ({
  value: selectorList = [],
  attributeList,
  onChange = noop,
}) => {
  const [valueOptionList, setValueOptionList] = useState<{key: string; value: string}[]>([]);
  const attributeOptionList = useMemo(
    () =>
      uniqBy(attributeList, 'key').map(({key}) => {
        return {
          key,
          value: key,
        };
      }),
    [attributeList]
  );

  const handleOnStepEntry = useCallback(
    (name, value) => {
      if (name === 'attribute') {
        const valueList = attributeList.reduce<{key: string; value: string}[]>((list, {key, value: attributeValue}) => {
          const itExists = list.find(item => item.value === attributeValue);
          return key === value && !itExists
            ? list.concat({
                key: attributeValue,
                value: attributeValue,
              })
            : list;
        }, []);
        setValueOptionList(valueList);
      }
    },
    [attributeList]
  );

  const onEntry = useCallback(
    (entry: string[]) => {
      const [attribute, operator, value] = entry;

      const selector: IItemSelector = {
        propertyName: attribute,
        value,
        operator: operator as CompareOperator,
        valueType: 'stringValue',
        locationName: LOCATION_NAME.SPAN_ATTRIBUTES,
      };

      onChange([...selectorList, selector]);
    },
    [onChange, selectorList]
  );

  const handleDeleteItemSelector = useCallback(
    (entryNumber: number) => {
      onChange(selectorList.filter((item, index) => index !== entryNumber));
    },
    [onChange, selectorList]
  );

  const handleClear = useCallback(() => {
    onChange([]);
  }, [onChange]);

  const defaultValueList = useMemo(
    () =>
      selectorList.flatMap(({propertyName, operator = CompareOperator.EQUALS, value}, index) => {
        return [
          {
            key: `${index}-${propertyName}`,
            label: propertyName,
            value: `${propertyName}${SEPARATOR}0${SEPARATOR}${index}`,
          },
          {
            key: `${index}-${operator}`,
            label: '=',
            value: `${operator}${SEPARATOR}1${SEPARATOR}${index}`,
          },
          {
            key: `${index}-${value}`,
            label: value,
            value: `${value}${SEPARATOR}2${SEPARATOR}${index}`,
          },
        ];
      }),
    [selectorList]
  );

  const entryList = useMemo(
    () => [
      {
        name: 'Attribute',
        items: attributeOptionList,
      },
      {
        name: 'Operator',
        items: operatorOptionList,
      },
      {
        name: 'Value',
        items: valueOptionList,
      },
    ],
    [attributeOptionList, valueOptionList]
  );

  return (
    <S.SelectorContainer data-cy="assertion-form-selector-input">
      <MultiSelectInput
        placeholder="Filter Spans"
        onClear={handleClear}
        entryList={entryList}
        onStepEntry={handleOnStepEntry}
        onEntry={onEntry}
        onDeselect={handleDeleteItemSelector}
        defaultValueList={defaultValueList}
      />
    </S.SelectorContainer>
  );
};

export default AssertionFormSelectorInput;
