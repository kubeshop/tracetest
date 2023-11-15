import * as S from 'components/RunDetailLayout/RunDetailLayout.styled';
import DemoSelector from 'components/DemoSelector';
import {TDraftTest} from 'types/Test.types';

interface IProps {
  demos: TDraftTest[];
}

const HeaderRight = ({demos}: IProps) => {
  return demos.length ? (
    <S.Section $justifyContent="end">
      <DemoSelector demos={demos} />
    </S.Section>
  ) : (
    <S.Section $justifyContent="" />
  );
};

export default HeaderRight;
