import {Tabs} from 'antd';
import {capitalize} from 'lodash';
import {useEffect, useMemo, useRef, useState} from 'react';

import AttributeList from 'components/AttributeList';
import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import SpanAttributeService from 'services/SpanAttribute.service';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import {getObjectIncludesText} from 'utils/Common';
import * as S from './SpanDetail.styled';

interface IProps {
  assertions?: TResultAssertions;
  attributeList: TSpanFlatAttribute[];
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  searchText?: string;
  type: SemanticGroupNames;
}

const Attributes = ({assertions, attributeList, onCreateTestSpec, searchText, type}: IProps) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [topPosition, setTopPosition] = useState(0);
  const sectionList = useMemo(
    () => SpanAttributeService.getSpanAttributeSectionsList(attributeList, type!),
    [attributeList, type]
  );

  useEffect(() => {
    setTopPosition(containerRef?.current?.offsetTop ?? 0);
  }, [attributeList]);

  return (
    <S.TabContainer $top={topPosition} ref={containerRef}>
      <Tabs
        data-cy="span-details-attributes"
        onChange={tabName => TraceAnalyticsService.onChangeTab(tabName)}
        size="small"
      >
        {sectionList.map(({section, attributeList: attrList}) => (
          <Tabs.TabPane
            tab={
              <span>
                {capitalize(section)} {getObjectIncludesText(attrList, searchText) && <S.Dot />}
              </span>
            }
            key={section}
          >
            <AttributeList
              assertions={assertions}
              attributeList={attrList}
              onCreateTestSpec={onCreateTestSpec}
              searchText={searchText}
            />
          </Tabs.TabPane>
        ))}
      </Tabs>
    </S.TabContainer>
  );
};

export default Attributes;
