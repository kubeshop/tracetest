import {FC, useState} from 'react';
import CreateAssertionModal from 'components/CreateAssertionModal';
import * as S from '../SpanDetail.styled';
import TraceAnalyticsService from '../../../services/Analytics/TraceAnalytics.service';
import {useAppSelector} from '../../../redux/hooks';
import AssertionSelectors from '../../../selectors/Assertion.selectors';
import {ISpanDetailProps} from '../SpanDetail';
import {AssetionsSpanComponent} from './AssetionsSpanComponent';
import {GenericHttpSpanHeader} from './GenericHttpSpanHeader';

const {onAddAssertionButtonClick} = TraceAnalyticsService;

const GenericHttpSpanDetail: FC<ISpanDetailProps> = ({testId, span, resultId}) => {
  const [openCreateAssertion, setOpenCreateAssertion] = useState(false);

  const assertionsResultList = useAppSelector(
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
      </S.DetailsContainer>
      {/*<Attributes spanAttributeList={span?.attributeList} />*/}

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

export default GenericHttpSpanDetail;
