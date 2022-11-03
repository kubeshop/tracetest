import {FormItemProps, Select, Form} from 'antd';
import {uniqBy} from 'lodash';
import {useMemo, useState} from 'react';
import SpanAttributeService from 'services/SpanAttribute.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import {OtelReference} from '../hooks/useGetOTELSemanticConventionAttributesInfo';
import * as S from './AttributeField.styled';
import {useDropDownRenderComponent} from './useDropDownRenderComponent';

interface IProps extends FormItemProps {
  attributeList: TSpanFlatAttribute[];
  reference: OtelReference;
}

const AttributeField = ({reference, attributeList, ...props}: IProps) => {
  const [hoveredKey, setHoveredKey] = useState<string | undefined>(undefined);
  const [newAttribute, setNewAttribute] = useState<string | undefined>(undefined);

  const filteredAttributedList = useMemo(() => {
    const list = uniqBy(attributeList, 'key');

    return newAttribute ? [{key: newAttribute, value: ''}, ...list] : list;
  }, [attributeList, newAttribute]);

  return (
    <Form.Item {...props}>
      <S.Select
        placeholder="Select Attribute"
        showSearch
        dropdownRender={useDropDownRenderComponent(reference, hoveredKey)}
        dropdownStyle={hoveredKey ? {minWidth: 550, maxWidth: 550} : undefined}
        filterOption={(search, option) =>
          SpanAttributeService.getItMatchesAttributeByKey(reference, option?.key || '', search)
        }
        onSearch={searchValue => setNewAttribute(searchValue)}
      >
        {filteredAttributedList.map(({key}) => (
          <Select.Option key={key} value={key}>
            <div onFocus={() => {}} onMouseOver={() => setHoveredKey(key)}>
              {key}
            </div>
          </Select.Option>
        ))}
      </S.Select>
    </Form.Item>
  );
};

export default AttributeField;
