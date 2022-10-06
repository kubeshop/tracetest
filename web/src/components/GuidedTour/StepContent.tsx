import {Button, Typography} from 'antd';
import React from 'react';
import {TooltipRenderProps} from 'react-joyride';
import {Container, Divider, Header, Title, TitleContainer, TitleText} from './StepContent.styled';

export const ReactJoyrideTooltip =
  (onComplete: () => void) =>
  ({continuous, index, step, backProps, primaryProps, tooltipProps, size, skipProps, isLastStep}: TooltipRenderProps) =>
    (
      <div {...tooltipProps} style={{width: 300, background: 'white'}}>
        <Header>
          <TitleContainer>
            <Title>{step.title}</Title>
            <TitleText>{` ${index + 1} of ${size}`}</TitleText>
          </TitleContainer>
        </Header>
        <Container>
          <Typography.Text>{step.content}</Typography.Text>
        </Container>
        <Divider />
        <div style={{padding: 16, display: 'flex', justifyContent: 'space-between'}}>
          <div>
            <Button {...skipProps} type="text">
              Skip
            </Button>
          </div>
          <div style={{display: 'flex', justifyContent: 'flex-end'}}>
            {index > 0 && (
              <Button {...backProps} type="text">
                Back
              </Button>
            )}
            {continuous && (
              <Button
                {...primaryProps}
                style={{marginLeft: 8}}
                type="text"
                onClick={e => {
                  primaryProps.onClick(e);
                  if (isLastStep) onComplete();
                }}
              >
                Next
              </Button>
            )}
          </div>
        </div>
      </div>
    );
