import {TEditTestForm} from '../EditTestModal';
import EditRequestDetailsHttp from './Http';

const EditRequestDetailsMap = {
  http: EditRequestDetailsHttp,
};

interface IProps {
  type: 'http';
  form: TEditTestForm;
}

const EditRequestDetails = ({type, form}: IProps) => {
  const Component = EditRequestDetailsMap[type];

  return <Component form={form} />;
};

export default EditRequestDetails;
