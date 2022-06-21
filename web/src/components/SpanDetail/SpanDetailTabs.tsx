import {Tabs} from 'antd';
import {capitalize} from 'lodash';
import React, {useMemo} from 'react';

import {getObjectIncludesText} from 'utils/Common';
import {useSpan} from 'providers/Span/Span.provider';
import AttributeList from 'components/AttributeList';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import SpanAttributeService from 'services/SpanAttribute.service';
import {ISpanDetailsComponentProps} from './SpanDetail';
import * as S from './SpanDetail.styled';

const SpanDetailTabs: React.FC<ISpanDetailsComponentProps> = ({
  span: {attributeList = [], type} = {},
  onCreateAssertion,
  assertions,
}) => {
  const {searchText} = useSpan();
  const sectionList = useMemo(
    () => SpanAttributeService.getSpanAttributeSectionsList(attributeList, type!),
    [attributeList, type]
  );

  return (
    <Tabs data-cy="span-details-attributes" onChange={tabName => TraceAnalyticsService.onChangeTab(tabName)}>
      {sectionList.map(({section, attributeList: attrList}) => (
        <Tabs.TabPane
          tab={
            <span>
              {capitalize(section)} {getObjectIncludesText(attrList, searchText) && <S.Dot />}
            </span>
          }
          key={section}
        >
          <AttributeList assertions={assertions} attributeList={attrList} onCreateAssertion={onCreateAssertion} />
        </Tabs.TabPane>
      ))}
    </Tabs>
  );
};

export default SpanDetailTabs;
