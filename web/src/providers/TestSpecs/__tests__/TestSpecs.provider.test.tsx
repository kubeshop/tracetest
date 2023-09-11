import {render} from 'test-utils';
import TestSpecs from '../TestSpecs.provider';

describe('TestDefinitionProvider', () => {
  it('should render with the proper values', () => {
    const {getByText} = render(
      <TestSpecs testId="testId" runId={1}>
        <div>
          <p>Hello</p>
        </div>
      </TestSpecs>
    );

    expect(getByText('Hello')).toBeInTheDocument();
  });
});
