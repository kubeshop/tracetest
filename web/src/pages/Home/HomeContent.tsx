import {useEffect, useState} from 'react';
import {delay} from 'lodash';
// import {Button} from 'antd';
// import {InfoCircleOutlined} from '@ant-design/icons';
import CreateTestModal from 'components/CreateTest';
import {Steps} from 'components/GuidedTour/homeStepList';
import GuidedTourService, {GuidedTours} from 'services/GuidedTourService';
import useGuidedTour from 'hooks/useGuidedTour';
import TestList from './TestList';
import * as S from './Home.styled';

const HomeContent: React.FC = () => {
  const [openCreateTestModal, setOpenCreateTestModal] = useState(false);

  const {setCurrentStep, setIsOpen, currentStep, isOpen: isGuidOpen} = useGuidedTour(GuidedTours.Home);

  useEffect(() => {
    if (currentStep > 0 && !openCreateTestModal && isGuidOpen) {
      setOpenCreateTestModal(true);
      setCurrentStep(2);
      delay(() => setCurrentStep(1), 0);
    }
  }, [currentStep, openCreateTestModal, setCurrentStep, isGuidOpen]);

  return (
    <S.Wrapper>
      <S.PageHeader>
        <S.TitleText>All Tests</S.TitleText>
        <S.ActionContainer>
          {/* <Button
            size="large"
            type="link"
            icon={<InfoCircleOutlined />}
            onClick={() => {
              setCurrentStep(0);
              setIsOpen(true);
            }}
          >
            Guided tour
          </Button> */}
          <S.CreateTestButton
            data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.CreateTest)}
            type="primary"
            size="large"
            onClick={() => {
              setOpenCreateTestModal(true);
              if (isGuidOpen) delay(() => setCurrentStep(currentStep + 1), 1);
            }}
          >
            Create Test
          </S.CreateTestButton>
        </S.ActionContainer>
      </S.PageHeader>
      <TestList />
      <CreateTestModal visible={openCreateTestModal} onClose={() => setOpenCreateTestModal(false)} />
    </S.Wrapper>
  );
};

export default HomeContent;
