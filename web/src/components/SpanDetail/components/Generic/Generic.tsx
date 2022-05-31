import {Tabs} from 'antd';
import React from 'react';

import AttributeList from 'components/AttributeList';
import {ISpanDetailsComponentProps} from 'components/SpanDetail/SpanDetail';
import * as S from 'components/SpanDetail/SpanDetail.styled';

const Generic: React.FC<ISpanDetailsComponentProps> = ({
  assertions,
  onCreateAssertion,
  span: {attributeList = []} = {},
}) => {
  return (
    <S.SpanTabs data-cy="span-details-attributes">
      <Tabs.TabPane tab="Attribute list" key="span-attribute-list">
        <AttributeList attributeList={attributeList} onCreateAssertion={onCreateAssertion} assertions={assertions} />
      </Tabs.TabPane>
    </S.SpanTabs>
  );
};

export default Generic;
