import {ParentSize} from '@visx/responsive';

import {IDiagramComponentProps} from 'components/Diagram/Diagram';

import Visualization from './Visualization';

const Timeline = (props: IDiagramComponentProps) => (
  <ParentSize parentSizeStyles={{height: 'auto', overflowY: 'scroll'}}>
    {({width}) => <Visualization {...props} width={width} />}
  </ParentSize>
);

export default Timeline;
