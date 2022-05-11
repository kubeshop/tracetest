import {Tabs} from 'antd';
import React from 'react';
import SpanAttributesTable from '../../../SpanAttributesTable/SpanAttributesTable';
import Assertions from '../../Assertions';
import {ISpanDetailProps} from '../../SpanDetail';
import * as S from '../../SpanDetail.styled';
import Request from './Request';
import Response from './Response';

const Http: React.FC<ISpanDetailProps> = ({
  span: {attributeList = []} = {},
  span,
  assertionsResultList = [],
  testId,
  resultId,
}) => {
  return (
    <S.SpanTabs data-cy="span-details-attributes">
      <Tabs.TabPane tab="Assertion" key="span-assertion">
        <Assertions span={span} assertionsResultList={assertionsResultList} testId={testId} resultId={resultId} />
      </Tabs.TabPane>
      <Tabs.TabPane tab="Request" key="span-request">
        <Request attributeList={attributeList} />
      </Tabs.TabPane>
      <Tabs.TabPane tab="Response" key="span-response">
        <Response attributeList={attributeList} />
      </Tabs.TabPane>
      <Tabs.TabPane tab="Attribute list" key="span-attribute-list">
        <SpanAttributesTable spanAttributesList={attributeList} />
      </Tabs.TabPane>
    </S.SpanTabs>
  );
};

export default Http;
