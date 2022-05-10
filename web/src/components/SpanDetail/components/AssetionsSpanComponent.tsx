import {Typography} from 'antd';
import {ISpan} from 'types/Span.types';
import SkeletonTable from 'components/SkeletonTable';
import AssertionsResultTable from 'components/AssertionsTable/AssertionsTable';
import * as S from 'components/SpanDetail/SpanDetail.styled';

interface IProps {
  span?: ISpan;
  testId?: string;
  resultId?: string;
  assertionsResultList: any[];
}

export const AssetionsSpanComponent = ({assertionsResultList, testId, span, resultId}: IProps): JSX.Element => (
  <SkeletonTable loading={!span}>
    {assertionsResultList.length ? (
      assertionsResultList
        .sort((a, b) => (a.assertion.assertionId > b.assertion.assertionId ? -1 : 1))
        .map(({assertion, assertionResultList}, index) => (
          <AssertionsResultTable
            key={assertion.assertionId}
            assertionResults={assertionResultList}
            sort={index + 1}
            assertion={assertion}
            span={span!}
            resultId={resultId!}
            testId={testId!}
          />
        ))
    ) : (
      <S.DetailsEmptyStateContainer>
        <S.DetailsTableEmptyStateIcon />
        <Typography.Text disabled>No Data</Typography.Text>
      </S.DetailsEmptyStateContainer>
    )}
  </SkeletonTable>
);
