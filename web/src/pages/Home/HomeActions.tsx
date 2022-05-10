import {Steps} from '../../components/GuidedTour/homeStepList';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import HomeAnalyticsService from '../../services/Analytics/HomeAnalytics.service';
import * as S from './Home.styled';

const {onCreateTestClick} = HomeAnalyticsService;

interface IHomeActionsProps {
  onCreateTest(): void;
}

const HomeActions: React.FC<IHomeActionsProps> = ({onCreateTest}) => {
  return (
    <S.ActionContainer>
      {/* <Button
            size="large"
            type="link"
            icon={<InfoCircleOutlined />}
            onClick={() => {
              setCurrentStep(0);
              setIsOpen(true);
              onGuidedTourClick();
            }}
          >
            Guided tour
          </Button> */}
      <S.CreateTestButton
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.CreateTest)}
        data-cy="create-test-button"
        type="primary"
        size="large"
        onClick={() => {
          onCreateTestClick();
          onCreateTest();
          // if (isGuidOpen) delay(() => setCurrentStep(currentStep + 1), 1);
        }}
      >
        Create Test
      </S.CreateTestButton>
    </S.ActionContainer>
  );
};

export default HomeActions;
