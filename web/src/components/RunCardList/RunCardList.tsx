import ResultCard from 'components/RunCard';

import {TTestRun} from 'types/TestRun.types';
import * as S from './RunCardList.styled';

interface IProps {
  resultList: TTestRun[];
  testId: string;
}

const RunCardList = ({resultList, testId}: IProps) => (
  <S.Container data-cy="run-card-list">
    {resultList.map(run => (
      <ResultCard key={run.id} run={run} testId={testId} />
    ))}
  </S.Container>
);

export default RunCardList;
