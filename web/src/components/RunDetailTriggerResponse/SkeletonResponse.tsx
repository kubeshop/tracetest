import * as S from './RunDetailTriggerResponse.styled';

const SkeletonResponse = () => {
  return (
    <S.LoadingResponseBody>
      <S.TextHolder $width={50} />
      <S.TextHolder $width={100} />
      <S.TextHolder $width={100} />
      <S.TextHolder $width={100} />
      <S.TextHolder $width={80} />
    </S.LoadingResponseBody>
  );
};

export default SkeletonResponse;
