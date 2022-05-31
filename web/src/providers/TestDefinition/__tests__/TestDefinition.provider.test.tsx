import {render} from '@testing-library/react';
import {BrowserRouter} from 'react-router-dom';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import TestDefinitionProvider from '../TestDefinition.provider';

describe('TestDefinitionProvider', () => {
  it('should render with the proper values', () => {
    const {getByText} = render(
      <BrowserRouter>
        <TestDefinitionProvider testId="testId" runId="runId">
          <div>
            <p>Hello</p>
          </div>
        </TestDefinitionProvider>
      </BrowserRouter>,
      {wrapper: ReduxWrapperProvider}
    );

    expect(getByText('Hello')).toBeInTheDocument();
  });
});
