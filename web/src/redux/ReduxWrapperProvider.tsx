import * as React from 'react';
import {Provider} from 'react-redux';
import {store} from './store';

export const ReduxWrapperProvider: React.FC<{children: JSX.Element}> = ({children}) => {
  return <Provider store={store}>{children}</Provider>;
};
