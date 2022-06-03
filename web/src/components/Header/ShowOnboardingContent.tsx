import {Button, Divider, Typography} from 'antd';
import {Dispatch, SetStateAction} from 'react';
import {AddAssertionButton} from '../TraceDrawer/TraceDrawer.styled';

export const ShowOnboardingContent = (
  onGuidedTourClick: () => void,
  setIsOpen: Dispatch<SetStateAction<Boolean>>,
  setVisible: Dispatch<SetStateAction<boolean>>
) => (
  <div>
    <div style={{padding: 16}}>
      <Typography.Text style={{margin: 0}}>Walk though Tracetest features</Typography.Text>
    </div>
    <Divider style={{margin: 0}} />
    <div style={{display: 'flex', justifyContent: 'flex-end', padding: 16}}>
      <Button style={{marginRight: 16}} ghost onClick={() => setVisible(false)} type="primary">
        No thanks
      </Button>
      <AddAssertionButton
        onClick={() => {
          setIsOpen(true);
          setVisible(false);
          onGuidedTourClick();
        }}
      >
        Show me around
      </AddAssertionButton>
    </div>
  </div>
);
