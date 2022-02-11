import ReactJson from 'react-json-view';

interface IProps {
  json: JSON;
}
const TraceData = ({json}: IProps) => {
  return <ReactJson src={json} collapsed />;
};

export default TraceData;
