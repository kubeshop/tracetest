import {TriggerTypes} from 'constants/Test.constants';
import CreateTestProvider from 'providers/CreateTest';
import {useParams} from 'react-router-dom';
import Content from './Content';

const CreateTest = () => {
  const {triggerType = ''} = useParams();

  return (
    <CreateTestProvider>
      <Content triggerType={triggerType as TriggerTypes} />
    </CreateTestProvider>
  );
};

export default CreateTest;
