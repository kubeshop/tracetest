import {render} from '@testing-library/react';
import {ReduxWrapperProvider} from '../../../redux/ReduxWrapperProvider';
import {TestingModels} from '../../../utils/TestingModels';
import TraceAssertionsTable from '../TraceAssertionsTable';

test('Layout', () => {
  // const result = render(
  //   <ReduxWrapperProvider>
  //     <TraceAssertionsTable assertionResult={TestingModels.assertionResult} onSpanSelected={jest.fn()} />
  //   </ReduxWrapperProvider>
  // );
  // expect(result.container).toMatchSnapshot();
});
