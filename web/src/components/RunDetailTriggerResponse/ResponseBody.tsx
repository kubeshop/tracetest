import SkeletonResponse from './SkeletonResponse';
import CodeBlock from '../CodeBlock';

interface IProps {
  body?: string;
}

const ResponseBody = ({body = ''}: IProps) => (!body ? <SkeletonResponse /> : <CodeBlock value={body} />);

export default ResponseBody;
