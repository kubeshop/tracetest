import CreateAssertionModal from 'components/CreateAssertionModal';
import React, {FC, useState} from 'react';
import {useAppSelector} from '../../../redux/hooks';
import AssertionSelectors from '../../../selectors/Assertion.selectors';
import TraceAnalyticsService from '../../../services/Analytics/TraceAnalytics.service';
import {IAssertion, ISpanAssertionResult} from '../../../types/Assertion.types';
import {ISpanDetailProps} from '../SpanDetail';
import * as S from '../SpanDetail.styled';
import {AssetionsSpanComponent} from './AssetionsSpanComponent';
import {GenericHttpSpanHeader} from './GenericHttpSpanHeader';
import {HTTPAttributesTabs} from './HTTPAttributesTabs';

const {onAddAssertionButtonClick} = TraceAnalyticsService;

export const GenericHttpSpanDetail: FC<ISpanDetailProps> = ({testId, span, resultId}) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);

  const assertionsResultList = useAppSelector<{assertion: IAssertion; assertionResultList: ISpanAssertionResult[]}[]>(
    AssertionSelectors.selectAssertionResultListBySpan(testId, resultId, span?.spanId)
  );
  return (
    <>
      <S.DetailsContainer>
        <GenericHttpSpanHeader
          title={`Details for ${span?.type?.toUpperCase()} ${span?.name}`}
          onClick={() => {
            onAddAssertionButtonClick();
            setOpenCreateAssertion(true);
          }}
        />
        <AssetionsSpanComponent
          resultId={resultId}
          testId={testId}
          span={span}
          assertionsResultList={assertionsResultList}
        />
        <HTTPAttributesTabs attributeList={span?.attributeList || []} />
        <br />
        <br />
      </S.DetailsContainer>

      {span && testId && (
        <CreateAssertionModal
          key={`KEY_${span?.spanId}`}
          testId={testId}
          span={span}
          resultId={resultId!}
          open={openCreateAssertion}
          onClose={() => setOpenCreateAssertion(false)}
        />
      )}
    </>
  );
};
