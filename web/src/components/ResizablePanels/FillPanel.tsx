import {PropsWithChildren} from 'react';
import * as Spaces from 'react-spaces';

const FillPanel: React.FC<PropsWithChildren<{}>> = ({children}) => <Spaces.Fill>{children}</Spaces.Fill>;

export default FillPanel;
