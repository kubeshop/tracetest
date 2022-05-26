import {render} from '@testing-library/react';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import TestDefinitionProvider from '../TestDefinition.provider';

describe('TestDefinitionProvider', () => {
  it('should render with the proper values', () => {
    const {getByText} = render(
      <TestDefinitionProvider testId="testId" runId="runId">
        <div>
          <p>Hello</p>
        </div>
      </TestDefinitionProvider>,
      {wrapper: ReduxWrapperProvider}
    );

    expect(getByText('Hello')).toBeInTheDocument();
  });
});
