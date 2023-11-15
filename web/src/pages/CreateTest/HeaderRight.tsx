import * as S from 'components/RunDetailLayout/RunDetailLayout.styled';
import DemoSelector from 'components/DemoSelector';
import {useCreateTest} from 'providers/CreateTest';
import {useSettingsValues} from '../../providers/SettingsValues/SettingsValues.provider';

const HeaderRight = () => {
  const {demoList} = useCreateTest();
  const {demos} = useSettingsValues();

  return demoList.length && demos.length ? (
    <S.Section $justifyContent="end">
      <DemoSelector />
    </S.Section>
  ) : (
    <S.Section $justifyContent="" />
  );
};

export default HeaderRight;
