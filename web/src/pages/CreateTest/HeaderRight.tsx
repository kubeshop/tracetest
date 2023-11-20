import DemoSelector from 'components/DemoSelector';
import * as S from 'components/RunDetailLayout/RunDetailLayout.styled';
import VariableSetSelector from 'components/VariableSetSelector';
import {TDraftTest} from 'types/Test.types';

interface IProps {
  demos: TDraftTest[];
}

const HeaderRight = ({demos}: IProps) => (
  <S.Section $justifyContent="flex-end">
    <VariableSetSelector />
    {demos.length && <DemoSelector demos={demos} />}
  </S.Section>
);

export default HeaderRight;
