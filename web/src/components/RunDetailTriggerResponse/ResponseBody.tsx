import SkeletonResponse from './SkeletonResponse';
import CodeBlock from '../CodeBlock';

interface IProps {
  body?: string;
  bodyMimeType?: string;
}

const ResponseBody = ({body = '', bodyMimeType = ''}: IProps) =>
  !body ? <SkeletonResponse /> : <CodeBlock value={body} mimeType={bodyMimeType} />;

export default ResponseBody;
