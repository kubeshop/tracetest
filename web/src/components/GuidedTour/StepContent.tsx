import {Button, Typography} from 'antd';
import React from 'react';
import {TooltipRenderProps} from 'react-joyride';
import * as S from './StepContent.styled';

const StepContent = ({
  continuous,
  index,
  step,
  backProps,
  primaryProps,
  tooltipProps,
  size,
  skipProps,
  isLastStep,
}: TooltipRenderProps) => (
  <S.Container {...tooltipProps} data-cy="onboarding-container">
    <S.Header>
      <S.Title>{step.title}</S.Title>
      <S.TitleText data-cy="onboarding-step">{` ${index + 1} of ${size}`}</S.TitleText>
    </S.Header>

    <S.Body>
      <Typography.Text>{step.content}</Typography.Text>
    </S.Body>

    <S.Footer>
      <div>
        {!isLastStep && (
          <Button {...skipProps} type="text">
            Skip
          </Button>
        )}
      </div>

      <div>
        {index > 0 && !isLastStep && (
          <Button {...backProps} type="text">
            Prev
          </Button>
        )}
        {continuous && (
          <Button {...primaryProps} type="link" data-cy="onboarding-next">
            {isLastStep ? 'Done' : 'Next'}
          </Button>
        )}
      </div>
    </S.Footer>
  </S.Container>
);

export default StepContent;
