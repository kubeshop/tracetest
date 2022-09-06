import {render} from 'test-utils';
import TestDefinitionProvider from '../TestSpecs.provider';

describe('TestDefinitionProvider', () => {
  it('should render with the proper values', () => {
    const {getByText} = render(
      <TestDefinitionProvider testId="testId" runId="runId">
        <div>
          <p>Hello</p>
        </div>
      </TestDefinitionProvider>
    );

    expect(getByText('Hello')).toBeInTheDocument();
  });
});
