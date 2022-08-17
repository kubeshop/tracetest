import RunDetailLayout from 'components/RunDetailLayout';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';

const Content = () => {
  const {test} = useTestDefinition();

  return test ? <RunDetailLayout test={test} /> : null;
};

export default Content;
