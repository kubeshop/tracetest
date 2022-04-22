import {useState} from 'react';
import {delay} from 'lodash';
import {Button} from 'antd';
import {InfoCircleOutlined} from '@ant-design/icons';
import CreateTestModal from 'components/CreateTest';
import {Steps} from 'components/GuidedTour/homeStepList';
import TestList from './TestList';
import * as S from './Home.styled';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTourService';
import useGuidedTour from '../../components/GuidedTour/useGuidedTour';

const HomeContent: React.FC = () => {
  const [openCreateTestModal, setOpenCreateTestModal] = useState(false);

  const {setCurrentStep, setIsOpen, currentStep} = useGuidedTour(GuidedTours.Home);

  return (
    <S.Wrapper>
      <S.PageHeader>
        <S.TitleText>All Tests</S.TitleText>
        <S.ActionContainer>
          <Button
            size="large"
            type="link"
            icon={<InfoCircleOutlined />}
            onClick={() => {
              setCurrentStep(0);
              setIsOpen(true);
            }}
          >
            Guided tour
          </Button>
          <S.CreateTestButton
            data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.CreateTest)}
            type="primary"
            size="large"
            onClick={() => {
              setOpenCreateTestModal(true);
              delay(() => setCurrentStep(currentStep + 1), 1);
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
