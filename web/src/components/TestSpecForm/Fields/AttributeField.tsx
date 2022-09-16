import {Select} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {uniqBy} from 'lodash';
import {ReactElement, useMemo, useState} from 'react';
import SpanAttributeService from 'services/SpanAttribute.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import {OtelReference} from '../hooks/useGetOTELSemanticConventionAttributesInfo';
import * as S from './AttributeField.styled';
import {useDropDownRenderComponent} from './useDropDownRenderComponent';

interface IProps {
  attributeList: TSpanFlatAttribute[];
  reference: OtelReference;
  field: Pick<FormListFieldData, never>;
  name: number;
}

export const AttributeField = ({field, name, reference, attributeList}: IProps): ReactElement => {
  const [hoveredKey, setHoveredKey] = useState<string | undefined>(undefined);
  const [newAttribute, setNewAttribute] = useState<string | undefined>(undefined);

  const filteredAttributedList = useMemo(() => {
    const list = uniqBy(attributeList, 'key');

    return newAttribute ? [{key: newAttribute, value: ''}, ...list] : list;
  }, [attributeList, newAttribute]);

  return (
    <S.FormItem
      {...field}
      name={[name, 'attribute']}
      rules={[{required: true, message: 'Attribute is required'}]}
      data-cy="assertion-check-attribute"
      id="assertion-check-attribute"
      style={{flexBasis: '30%', width: 0}}
    >
      <S.Select
        placeholder="Select Attribute"
        showSearch
        dropdownRender={useDropDownRenderComponent(reference, hoveredKey)}
        dropdownStyle={hoveredKey ? {minWidth: 550, maxWidth: 550} : undefined}
        filterOption={(search, option) => {
          const itMatches = SpanAttributeService.getItMatchesAttributeByKey(reference, option?.key || '', search);

          return itMatches;
        }}
        onSearch={value => setNewAttribute(value)}
      >
        {filteredAttributedList.map(({key}) => (
          <Select.Option key={key} value={key}>
            <div onFocus={() => {}} onMouseOver={() => setHoveredKey(key)}>
              {key}
            </div>
          </Select.Option>
        ))}
      </S.Select>
    </S.FormItem>
  );
};
