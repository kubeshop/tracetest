import {Tabs} from 'antd';
import {capitalize} from 'lodash';
import React, {useMemo} from 'react';
import SpanAttributeService from 'services/SpanAttribute.service';
import AttributeList from '../AttributeList';
import {ISpanDetailsComponentProps} from './SpanDetail';
import * as S from './SpanDetail.styled';

const SpanDetailTabs: React.FC<ISpanDetailsComponentProps> = ({span: {attributeList = [], type} = {}, onCreateAssertion, assertions}) => {
  const sectionList = useMemo(
    () => SpanAttributeService.getSpanAttributeSectionsList(attributeList, type!),
    [attributeList, type]
  );

  return (
    <S.SpanTabs data-cy="span-details-attributes">
      {sectionList.map(({section, attributeList: attrList}) => (
        <Tabs.TabPane tab={capitalize(section)} key={section}>
          <AttributeList assertions={assertions} attributeList={attrList} onCreateAssertion={onCreateAssertion} />
        </Tabs.TabPane>
      ))}
    </S.SpanTabs>
  );
};

export default SpanDetailTabs;
