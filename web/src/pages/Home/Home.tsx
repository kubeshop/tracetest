import Title from 'antd/lib/typography/Title';
import {useState} from 'react';
import CreateTestModal from '../../components/CreateTestModal';
import TestList from './TestList';
import * as S from './Home.styled';

const Home = (): JSX.Element => {
  const [openCreateTestModal, setOpenCreateTestModal] = useState<boolean>(false);
  const handleCreateTest = () => {
    setOpenCreateTestModal(true);
  };

  return (
    <div>
      <S.Header>
        <Title level={2}>Tracetest</Title>
      </S.Header>
      <S.Content>
        <S.SideMenu>
          <S.CreateTestButton type="primary" size="large" onClick={handleCreateTest}>
            Create Test
          </S.CreateTestButton>
        </S.SideMenu>
        <S.TestsContainer>
          <TestList />
        </S.TestsContainer>
      </S.Content>
      <CreateTestModal visible={openCreateTestModal} onClose={() => setOpenCreateTestModal(false)} />
    </div>
  );
};

export default Home;
