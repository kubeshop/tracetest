import {THeader} from 'types/Test.types';
import HeaderRow from 'components/HeaderRow';
import SkeletonResponse from './SkeletonResponse';
import * as S from './RunDetailTriggerResponse.styled';

interface IProps {
  headers?: THeader[];
}

const ResponseHeaders = ({headers}: IProps) => {
  return !headers ? (
    <SkeletonResponse />
  ) : (
    <S.HeadersList>
      {headers.map(header => (
        <HeaderRow header={header} key={header.key} />
      ))}
    </S.HeadersList>
  );
};

export default ResponseHeaders;
