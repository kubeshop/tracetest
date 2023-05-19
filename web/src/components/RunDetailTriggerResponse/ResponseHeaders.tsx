import HeaderRow from 'components/HeaderRow';
import {isRunStateFinished} from 'models/TestRun.model';
import {THeader} from 'types/Test.types';
import {TTestRunState} from 'types/TestRun.types';
import SkeletonResponse from './SkeletonResponse';
import * as S from './RunDetailTriggerResponse.styled';

interface IProps {
  headers?: THeader[];
  state: TTestRunState;
}

const ResponseHeaders = ({headers, state}: IProps) =>
  isRunStateFinished(state) || !!headers ? (
    <S.HeadersList>{headers && headers.map(header => <HeaderRow header={header} key={header.key} />)}</S.HeadersList>
  ) : (
    <SkeletonResponse />
  );

export default ResponseHeaders;
