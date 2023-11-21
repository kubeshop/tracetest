import {TriggerTypes} from 'constants/Test.constants';
import {useParams} from 'react-router-dom';
import Content from './Content';

const CreateTest = () => {
  const {triggerType = ''} = useParams();

  return <Content triggerType={triggerType as TriggerTypes} />;
};

export default CreateTest;
