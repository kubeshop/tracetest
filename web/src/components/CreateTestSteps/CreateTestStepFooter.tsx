import {Button} from 'antd';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import * as S from './CreateTestSteps.styled';

interface IProps {
  isValid: boolean;
  onNext(): void;
}

const CreateStepFooter = ({isValid, onNext}: IProps) => {
  const {stepNumber, stepList, onPrev, isLoading} = useCreateTest();

  return (
    <S.Footer>
      {Boolean(stepNumber) && (
        <Button data-cy="create-test-prev-button" type="text" onClick={onPrev}>
          Previous
        </Button>
      )}
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
    </S.Footer>
  );
};

export default CreateStepFooter;
