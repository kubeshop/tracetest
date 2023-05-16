import {isRunStateFinished} from 'models/TestRun.model';
import {TTestRunState} from 'types/TestRun.types';
import SkeletonResponse from './SkeletonResponse';
import CodeBlock from '../CodeBlock';

interface IProps {
  body?: string;
  bodyMimeType?: string;
  state: TTestRunState;
}

const ResponseBody = ({body = '', bodyMimeType = '', state}: IProps) =>
  isRunStateFinished(state) || !!body ? <CodeBlock value={body} mimeType={bodyMimeType} /> : <SkeletonResponse />;

export default ResponseBody;
