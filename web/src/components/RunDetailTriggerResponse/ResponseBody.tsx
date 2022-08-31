import SkeletonResponse from './SkeletonResponse';
import * as S from './RunDetailTriggerResponse.styled';

interface IProps {
  body?: string;
}

const ResponseBody = ({body = ''}: IProps) => {
  return !body ? (
    <SkeletonResponse />
  ) : (
    <S.ValueJson>
      <pre>{JSON.stringify(JSON.parse(body), null, 2)}</pre>
    </S.ValueJson>
  );
};

export default ResponseBody;
