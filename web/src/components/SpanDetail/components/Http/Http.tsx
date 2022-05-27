import {Tabs} from 'antd';
import React, {useMemo} from 'react';

import {HttpRequestAttributeList, HttpResponseAttributeList} from 'constants/Span.constants';
import AttributeList from 'components/AttributeList';
import {ISpanDetailsComponentProps} from 'components/SpanDetail/SpanDetail';
import * as S from 'components/SpanDetail/SpanDetail.styled';
import {TSpanFlatAttribute} from 'types/Span.types';

const filterRequestList = (attributeList: TSpanFlatAttribute[]) =>
  attributeList?.filter(a => HttpRequestAttributeList.includes(a.key) || a.key.includes('http.request'));

const filterResponseList = (attributeList: TSpanFlatAttribute[]) =>
  attributeList?.filter(a => HttpResponseAttributeList.includes(a.key) || a.key.includes('http.response'));

const Http: React.FC<ISpanDetailsComponentProps> = ({
  assertions,
  onCreateAssertion,
  span: {attributeList = []} = {},
}) => {
  const requestList = useMemo(() => filterRequestList(attributeList), [attributeList]);
  const responseList = useMemo(() => filterResponseList(attributeList), [attributeList]);

  return (
    <S.SpanTabs data-cy="span-details-attributes">
      <Tabs.TabPane tab="Request" key="span-request">
        <AttributeList attributeList={requestList} onCreateAssertion={onCreateAssertion} assertions={assertions} />
      </Tabs.TabPane>
      <Tabs.TabPane tab="Response" key="span-response">
        <AttributeList attributeList={responseList} onCreateAssertion={onCreateAssertion} assertions={assertions} />
      </Tabs.TabPane>
      <Tabs.TabPane tab="Attribute list" key="span-attribute-list">
        <AttributeList attributeList={attributeList} onCreateAssertion={onCreateAssertion} assertions={assertions} />
      </Tabs.TabPane>
    </S.SpanTabs>
  );
};

export default Http;
