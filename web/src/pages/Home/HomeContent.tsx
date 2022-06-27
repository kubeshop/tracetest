import TestList from './TestList';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import SearchInput from '../../components/SearchInput';

const HomeContent: React.FC = () => {
  return (
    <S.Wrapper>
      <S.TitleText level={4}>All Tests</S.TitleText>
      <S.PageHeader>
        <SearchInput onSearch={() => console.log('onSearch')} placeholder="Search test (Not implemented yet)" />
        <HomeActions />
      </S.PageHeader>
      <TestList />
    </S.Wrapper>
  );
};

export default HomeContent;
