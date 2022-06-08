import {BtnFnProps} from '@reactour/tour/dist/types';
import {Typography} from 'antd';

export const PreviousButton: React.FC<BtnFnProps> = ({currentStep, stepsLength}) => (
  <div>
    <Typography.Text>{` ${currentStep + 1} of ${stepsLength}`}</Typography.Text>
  </div>
);
