import {useState} from 'react';
import SearchInput from '../../components/SearchInput';
import useInfiniteScroll from '../../hooks/useInfiniteScroll';
import {useGetTestListQuery} from '../../redux/apis/TraceTest.api';
import {TTest} from '../../types/Test.types';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import TestList from './TestList';

const HomeContent: React.FC = () => {
  const [state, setState] = useState<{query: string}>(() => ({query: '', filteredList: []}));
  const query = useInfiniteScroll<TTest, {query: string}>(useGetTestListQuery, {query: state.query});
  return (
    <S.Wrapper>
      <S.TitleText>All Tests</S.TitleText>
      <S.PageHeader>
        <SearchInput onSearch={value => setState(st => ({...st, query: value}))} placeholder="Search test" />
        <HomeActions />
      </S.PageHeader>
      <TestList query={query} />
    </S.Wrapper>
  );
};

export default HomeContent;
