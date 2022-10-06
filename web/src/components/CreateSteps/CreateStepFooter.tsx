import {Button} from 'antd';
import {ICreateTestStep} from 'types/Plugins.types';
import * as S from './CreateSteps.styled';

interface IProps {
  isValid: boolean;
  stepNumber: number;
  isLastStep: boolean;
  step: ICreateTestStep;
  onPrev(): void;
  isLoading: boolean;
}

const CreateStepFooter = ({isValid, stepNumber, step, isLastStep, onPrev, isLoading}: IProps) => {
  return (
    <S.Footer>
      <span>
        {stepNumber > 0 && (
          <Button data-cy="create-test-prev-button" type="primary" ghost onClick={onPrev}>
            Previous
          </Button>
        )}
      </span>
      <span>
        {!isLastStep ? (
          <Button
            htmlType="submit"
            form={step.component}
            data-cy="create-test-next-button"
            disabled={!isValid}
            type="primary"
          >
            Next
          </Button>
        ) : (
          <Button
            htmlType="submit"
            form={step.component}
            data-cy="create-test-create-button"
            disabled={!isValid}
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
