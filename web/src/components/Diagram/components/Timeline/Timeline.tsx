import {ParentSize} from '@visx/responsive';

import {IDiagramComponentProps} from 'components/Diagram/Diagram';
import Controls from '../DAG/Controls';

import Visualization from './Visualization';

const Timeline = (props: IDiagramComponentProps) => (
  <>
    <Controls mode="timeline" />
    <ParentSize parentSizeStyles={{height: 'auto', overflowY: 'scroll', marginTop: 16}}>
      {({width}) => <Visualization {...props} width={width} />}
    </ParentSize>
  </>
);

export default Timeline;
