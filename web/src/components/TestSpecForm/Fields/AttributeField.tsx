import {Select} from 'antd';
import {FormListFieldData} from 'antd/lib/form/FormList';
import {uniqBy} from 'lodash';
import React, {ReactElement, useMemo, useState} from 'react';
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
  return (
    <S.FormItem
      {...field}
      name={[name, 'attribute']}
      rules={[{required: true, message: 'Attribute is required'}]}
      data-cy="assertion-check-attribute"
      id="assertion-check-attribute"
    >
      <S.Select
        placeholder="Attribute"
        showSearch
        dropdownRender={useDropDownRenderComponent(reference, hoveredKey)}
        dropdownStyle={hoveredKey ? {minWidth: 550, maxWidth: 550} : undefined}
        filterOption={(input, option) => {
          const key = option?.key || '';
          const ref = SpanAttributeService.referencePicker(reference, key);
          const availableTagsMatchInput = Boolean(
            ref.tags.find(tag => tag.toString().toLowerCase().includes(input.toLowerCase()))
          );
          const currentOptionMatchInput = option?.key.includes(input);
          const currentDescriptionMatchInput = ref?.description.includes(input);
          return availableTagsMatchInput || currentOptionMatchInput || currentDescriptionMatchInput;
        }}
      >
        {useMemo(() => {
          return uniqBy(attributeList, 'key').map(({key}) => (
            <Select.Option key={key} value={key}>
              <div onFocus={() => {}} onMouseOver={() => setHoveredKey(key)}>
                {key}
              </div>
            </Select.Option>
          ));
        }, [attributeList])}
      </S.Select>
    </S.FormItem>
  );
};
