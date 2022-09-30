import {useState} from 'react';

import SearchInput from 'components/SearchInput';
import CreateTestModal from 'components/CreateTestModal/CreateTestModal';
import CreateTransactionModal from 'components/CreateTransactionModal/CreateTransactionModal';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import TestList from './TestList';

const Content = () => {
  const [query, setQuery] = useState<string>('');
  const [isCreateTransactionOpen, setIsCreateTransactionOpen] = useState(false);
  const [isCreateTestOpen, setIsCreateTestOpen] = useState(false);

  return (
    <>
      <S.Wrapper>
        <S.HeaderContainer>
          <S.TitleText>All Tests</S.TitleText>
        </S.HeaderContainer>
        <S.PageHeader>
          <SearchInput onSearch={value => setQuery(value)} placeholder="Search test" />
          <HomeActions
            onCreateTest={() => setIsCreateTestOpen(true)}
            onCreateTransaction={() => setIsCreateTransactionOpen(true)}
          />
        </S.PageHeader>
        <TestList query={query} />
      </S.Wrapper>
      <CreateTestModal isOpen={isCreateTestOpen} onClose={() => setIsCreateTestOpen(false)} />
      <CreateTransactionModal isOpen={isCreateTransactionOpen} onClose={() => setIsCreateTransactionOpen(false)} />
    </>
  );
};

export default Content;
