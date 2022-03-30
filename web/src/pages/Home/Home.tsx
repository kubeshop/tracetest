import {useState} from 'react';
import CreateTestModal from 'components/CreateTestModal';
import TestList from './TestList';
import * as S from './Home.styled';
import Layout from '../../components/Layout';

const Home = (): JSX.Element => {
  const [openCreateTestModal, setOpenCreateTestModal] = useState<boolean>(false);
  const handleCreateTest = () => {
    setOpenCreateTestModal(true);
  };

  return (
    <Layout>
      <S.Wrapper>
        <S.PageHeader>
          <S.TitleText>All Tests</S.TitleText>
          <S.CreateTestButton type="primary" size="large" onClick={handleCreateTest}>
            Create Test
          </S.CreateTestButton>
        </S.PageHeader>
        <TestList />
        <CreateTestModal visible={openCreateTestModal} onClose={() => setOpenCreateTestModal(false)} />
      </S.Wrapper>
    </Layout>
  );
};

export default Home;
