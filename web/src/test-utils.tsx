import {FC, PropsWithChildren, ReactElement} from 'react';
import {render, RenderOptions} from '@testing-library/react';
import {MemoryRouter} from 'react-router-dom';
import {ThemeProvider} from 'styled-components';

import {theme} from 'constants/Theme.constants';
import {ReduxWrapperProvider} from 'redux/ReduxWrapperProvider';

const Providers: FC<PropsWithChildren<{}>> = ({children}) => {
  return (
    <ThemeProvider theme={theme}>
      <ReduxWrapperProvider>
        <MemoryRouter>{children}</MemoryRouter>
      </ReduxWrapperProvider>
    </ThemeProvider>
  );
};

const customRender = (ui: ReactElement, options?: Omit<RenderOptions, 'wrapper'>) =>
  render(ui, {wrapper: Providers, ...options});

// re-export everything
export * from '@testing-library/react';
// override render method
export {customRender as render};
