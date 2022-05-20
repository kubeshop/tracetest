import {Tabs} from 'antd';
import React, {useMemo} from 'react';
import {HttpRequestAttributeList, HttpResponseAttributeList} from '../../../../constants/Span.constants';
import {TSpanFlatAttribute} from '../../../../types/Span.types';
import AttributeList from '../../../AttributeList';
import {TSpanDetailsComponentProps} from '../../SpanDetail';
import * as S from '../../SpanDetail.styled';

const filterRequestList = (attributeList: TSpanFlatAttribute[]) =>
  attributeList?.filter(a => HttpRequestAttributeList.includes(a.key) || a.key.includes('http.request'));

const filterResponseList = (attributeList: TSpanFlatAttribute[]) =>
  attributeList?.filter(a => HttpResponseAttributeList.includes(a.key) || a.key.includes('http.request'));

const Http: React.FC<TSpanDetailsComponentProps> = ({span: {attributeList = []} = {}, onCreateAssertion}) => {
  const responseList = useMemo(() => filterResponseList(attributeList), [attributeList]);
  const requestList = useMemo(() => filterRequestList(attributeList), [attributeList]);

  return (
    <S.SpanTabs data-cy="span-details-attributes">
      <Tabs.TabPane tab="Request" key="span-request">
        <AttributeList attributeList={requestList} onCreateAssertion={onCreateAssertion} />
      </Tabs.TabPane>
      <Tabs.TabPane tab="Response" key="span-response">
        <AttributeList attributeList={responseList} onCreateAssertion={onCreateAssertion} />
      </Tabs.TabPane>
      <Tabs.TabPane tab="Attribute list" key="span-attribute-list">
        <AttributeList attributeList={attributeList} onCreateAssertion={onCreateAssertion} />
      </Tabs.TabPane>
    </S.SpanTabs>
  );
};

export default Http;
