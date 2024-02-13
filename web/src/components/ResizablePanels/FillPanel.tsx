import * as Spaces from 'react-spaces';

const FillPanel: React.FC = ({children}) => <Spaces.Fill style={{overflow: 'scroll'}}>{children}</Spaces.Fill>;

export default FillPanel;
