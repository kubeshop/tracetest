import {FC, ReactElement} from 'react';
import {render, RenderOptions} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';

import {ReduxWrapperProvider} from 'redux/ReduxWrapperProvider';

const Providers: FC = ({children}) => {
  return (
    <ReduxWrapperProvider>
      <MemoryRouter>{children}</MemoryRouter>
    </ReduxWrapperProvider>
  );
};

const customRender = (ui: ReactElement, options?: Omit<RenderOptions, 'wrapper'>) =>
  render(ui, {wrapper: Providers, ...options});

// re-export everything
export * from '@testing-library/react';
// override render method
export {customRender as render};
