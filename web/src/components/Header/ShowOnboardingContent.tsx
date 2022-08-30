import {Button, Divider, Typography} from 'antd';
import {PrimaryButton} from '../TestResults/TestResults.styled';

export const ShowOnboardingContent = (onGuidedTourClick: () => void, onShow: () => void, onClose: () => void) => (
  <div>
    <div style={{padding: 16}}>
      <Typography.Text style={{margin: 0}}>Walk though Tracetest features</Typography.Text>
    </div>
    <Divider style={{margin: 0}} />
    <div style={{display: 'flex', justifyContent: 'flex-end', padding: 16}}>
      <Button data-cy="no-thanks" style={{marginRight: 16}} ghost onClick={() => onClose()} type="primary">
        No thanks
      </Button>
      <PrimaryButton
        onClick={() => {
          onShow();
          onClose();
          onGuidedTourClick();
        }}
      >
        Show me around
      </PrimaryButton>
    </div>
  </div>
);
