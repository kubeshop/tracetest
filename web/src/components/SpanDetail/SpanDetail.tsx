import {SemanticGroupNames} from '../../constants/SemanticGroupNames.constants';
import {ISpan} from '../../types/Span.types';
import GenericSpanDetail from './components/GenericSpanDetail';
import {GenericHttpSpanDetail} from './components/GenericHttpSpanDetail';

export interface ISpanDetailProps {
  testId?: string;
  span?: ISpan;
  resultId?: string;
}

const ComponentMap: Record<string, typeof GenericSpanDetail> = {
  [SemanticGroupNames.Http]: GenericHttpSpanDetail,
};

const SpanDetail: React.FC<ISpanDetailProps> = ({span, ...props}) => {
  const Component = ComponentMap[span?.type || ''] || GenericSpanDetail;

  return <Component span={span} {...props} />;
};

export default SpanDetail;
