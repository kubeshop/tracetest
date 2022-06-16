import {useEffect, useState} from 'react';
import {delay} from 'lodash';
import TestList from './TestList';
import * as S from './Home.styled';
import HomeActions from './HomeActions';
import CreateTestModal from '../../components/CreateTestModal';
import SearchInput from '../../components/SearchInput';
import {useGuidedTour} from '../../providers/GuidedTour/GuidedTour.provider';

const HomeContent: React.FC = () => {
  const [openCreateTestModal, setOpenCreateTestModal] = useState(false);
  const {
    tour: {setCurrentStep, currentStep, isOpen},
  } = useGuidedTour();

  useEffect(() => {
    if (currentStep > 0 && !openCreateTestModal && isOpen) {
      setOpenCreateTestModal(true);
      setCurrentStep(2);
      delay(() => setCurrentStep(1), 0);
    }
  }, [currentStep, isOpen, openCreateTestModal, setCurrentStep]);

  return (
    <>
      <S.Wrapper>
        <S.TitleText level={4}>All Tests</S.TitleText>
        <S.PageHeader>
          <SearchInput onSearch={() => console.log('onSearch')} placeholder="Search test (Not implemented yet)" />
          <HomeActions onCreateTest={() => setOpenCreateTestModal(true)} />
        </S.PageHeader>
        <TestList />
      </S.Wrapper>
      <CreateTestModal visible={openCreateTestModal} onClose={() => setOpenCreateTestModal(false)} />
    </>
  );
};

export default HomeContent;
