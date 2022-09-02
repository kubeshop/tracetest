import RunDetailLayout from 'components/RunDetailLayout';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';

const Content = () => {
  const {test} = useTestSpecs();

  return test ? <RunDetailLayout test={test} /> : null;
};

export default Content;
