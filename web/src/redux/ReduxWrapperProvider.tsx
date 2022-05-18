// test-utils.jsx
import * as React from 'react';
import {Provider} from 'react-redux';
import {store} from './store';

export const ReduxWrapperProvider: React.FC = ({children}) => {
  return <Provider store={store}>{children}</Provider>;
};
