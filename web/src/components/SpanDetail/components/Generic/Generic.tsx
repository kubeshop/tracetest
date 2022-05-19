import {Tabs} from 'antd';
import React from 'react';
import AttributeList from '../../../AttributeList';
import {ISpanDetailsComponentProps} from '../../SpanDetail';
import * as S from '../../SpanDetail.styled';

const Generic: React.FC<ISpanDetailsComponentProps> = ({span: {attributeList = []} = {}, onCreateAssertion}) => {
  return (
    <S.SpanTabs data-cy="span-details-attributes">
      <Tabs.TabPane tab="Attribute list" key="span-attribute-list">
        <AttributeList attributeList={attributeList} onCreateAssertion={onCreateAssertion} />
      </Tabs.TabPane>
    </S.SpanTabs>
  );
};

export default Generic;
