import {THeader} from 'types/Test.types';
import HeaderRow from './HeaderRow';
import SkeletonResponse from './SkeletonResponse';

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
    <>
      {headers.map(header => (
        <HeaderRow onCopy={onCopy} header={header} key={header.key} />
      ))}
    </>
  );
};

export default ResponseHeaders;
