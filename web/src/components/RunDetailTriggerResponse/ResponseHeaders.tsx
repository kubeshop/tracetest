import {THeader} from 'types/Test.types';
import HeaderRow from 'components/HeaderRow';
import SkeletonResponse from './SkeletonResponse';
import * as S from './RunDetailTriggerResponse.styled';

interface IProps {
  headers?: THeader[];
}

const ResponseHeaders = ({headers}: IProps) => {
  const onCopy = (value: string) => {
    navigator.clipboard.writeText(value);
  };

  return !headers ? (
    <SkeletonResponse />
  ) : (
    <S.HeadersList>
      {headers.map(header => (
        <HeaderRow onCopy={onCopy} header={header} key={header.key} />
      ))}
    </S.HeadersList>
  );
};

export default ResponseHeaders;
