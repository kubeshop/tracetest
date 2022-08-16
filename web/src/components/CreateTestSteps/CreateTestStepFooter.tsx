import {Button} from 'antd';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {useNavigate} from 'react-router-dom';
import * as S from './CreateTestSteps.styled';

interface IProps {
  isValid: boolean;
  onNext(): void;
}

const CreateStepFooter = ({isValid, onNext}: IProps) => {
  const {stepNumber, stepList, onPrev, isLoading} = useCreateTest();

  const navigateFunction = useNavigate();
  return (
    <S.Footer>
      {stepNumber > 0 ? (
        <Button data-cy="create-test-prev-button" type="primary" ghost onClick={onPrev}>
          Previous
        </Button>
      ) : (
        <div />
      )}
      <span>
        <Button
          style={{marginRight: 16}}
          data-cy="create-test-cancel"
          onClick={() => navigateFunction('/')}
          type="primary"
          ghost
        >
          Cancel
        </Button>
        {stepNumber < stepList.length - 1 ? (
          <Button data-cy="create-test-next-button" disabled={!isValid} onClick={onNext} type="primary">
            Next
          </Button>
        ) : (
          <Button
            data-cy="create-test-create-button"
            disabled={!isValid}
            onClick={onNext}
            type="primary"
            loading={isLoading}
          >
            Create
          </Button>
        )}
      </span>
    </S.Footer>
  );
};

export default CreateStepFooter;
