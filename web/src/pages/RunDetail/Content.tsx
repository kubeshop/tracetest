import RunDetailLayout from 'components/RunDetailLayout';
import {useTest} from 'providers/Test/Test.provider';

const Content = () => {
  const {test} = useTest();

  return test ? <RunDetailLayout test={test} /> : null;
};

export default Content;
