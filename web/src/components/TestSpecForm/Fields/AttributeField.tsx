import {Select} from 'antd';
import {noop, uniqBy} from 'lodash';
import {ReactElement, useMemo, useState} from 'react';
import SpanAttributeService from 'services/SpanAttribute.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import {OtelReference} from '../hooks/useGetOTELSemanticConventionAttributesInfo';
import * as S from './AttributeField.styled';
import {useDropDownRenderComponent} from './useDropDownRenderComponent';

interface IProps {
  attributeList: TSpanFlatAttribute[];
  reference: OtelReference;
  value?: string;
  onChange?(value: string): void;
}

export const AttributeField = ({reference, attributeList, value = '', onChange = noop}: IProps): ReactElement => {
  const [hoveredKey, setHoveredKey] = useState<string | undefined>(undefined);
  const [newAttribute, setNewAttribute] = useState<string | undefined>(undefined);

  const filteredAttributedList = useMemo(() => {
    const list = uniqBy(attributeList, 'key');

    return newAttribute ? [{key: newAttribute, value: ''}, ...list] : list;
  }, [attributeList, newAttribute]);

  return (
    <S.Select
      placeholder="Select Attribute"
      showSearch
      value={value}
      onChange={newValue => onChange(newValue as string)}
      dropdownRender={useDropDownRenderComponent(reference, hoveredKey)}
      dropdownStyle={hoveredKey ? {minWidth: 550, maxWidth: 550} : undefined}
      filterOption={(search, option) => {
        const itMatches = SpanAttributeService.getItMatchesAttributeByKey(reference, option?.key || '', search);

        return itMatches;
      }}
      onSearch={sarchValue => setNewAttribute(sarchValue)}
    >
      {filteredAttributedList.map(({key}) => (
        <Select.Option key={key} value={key}>
          <div onFocus={() => {}} onMouseOver={() => setHoveredKey(key)}>
            {key}
          </div>
        </Select.Option>
      ))}
    </S.Select>
  );
};
