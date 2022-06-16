import {Tabs} from 'antd';
import {capitalize} from 'lodash';
import React, {useMemo} from 'react';

import AttributeList from 'components/AttributeList';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import SpanAttributeService from 'services/SpanAttribute.service';
import {ISpanDetailsComponentProps} from './SpanDetail';

const SpanDetailTabs: React.FC<ISpanDetailsComponentProps> = ({
  span: {attributeList = [], type} = {},
  onCreateAssertion,
  assertions,
}) => {
  const sectionList = useMemo(
    () => SpanAttributeService.getSpanAttributeSectionsList(attributeList, type!),
    [attributeList, type]
  );

  return (
    <Tabs data-cy="span-details-attributes" onChange={tabName => TraceAnalyticsService.onChangeTab(tabName)}>
      {sectionList.map(({section, attributeList: attrList}) => (
        <Tabs.TabPane tab={capitalize(section)} key={section}>
          <AttributeList assertions={assertions} attributeList={attrList} onCreateAssertion={onCreateAssertion} />
        </Tabs.TabPane>
      ))}
    </Tabs>
  );
};

export default SpanDetailTabs;
