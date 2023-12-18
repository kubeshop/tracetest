import {Button} from 'antd';
import {ICreateTestStep} from 'types/Plugins.types';
import CreateTestSuiteAnalytics from 'services/Analytics/CreateTestSuite.service';
import * as S from './CreateSteps.styled';

interface IProps {
  isValid: boolean;
  stepNumber: number;
  isLastStep: boolean;
  step: ICreateTestStep;
  onPrev(): void;
  isLoading: boolean;
  mode: string;
}

const CreateStepFooter = ({isValid, stepNumber, step, isLastStep, onPrev, isLoading, mode}: IProps) => {
  return (
    <S.Footer>
      <span>
        {stepNumber > 0 && (
          <Button
            data-cy="create-test-prev-button"
            type="primary"
            ghost
            onClick={() => {
              CreateTestSuiteAnalytics.onPrevClick(step.name);
              onPrev();
            }}
          >
            Previous
          </Button>
        )}
      </span>
      <span>
        {!isLastStep ? (
          <Button
            htmlType="submit"
            form={step.component}
            data-cy={`${mode}-create-next-button`}
            disabled={!isValid}
            type="primary"
            onClick={() => CreateTestSuiteAnalytics.onNextClick(step.name)}
          >
            Next
          </Button>
        ) : (
          <Button
            htmlType="submit"
            form={step.component}
            onClick={() => CreateTestSuiteAnalytics.onFormSubmit()}
            data-cy={`${mode}-create-create-button`}
            disabled={!isValid}
            type="primary"
            loading={isLoading}
          >
            Create & Run
          </Button>
        )}
      </span>
    </S.Footer>
  );
};

export default CreateStepFooter;
