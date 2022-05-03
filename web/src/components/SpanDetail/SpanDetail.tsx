import {SemanticGroupNames} from '../../constants/SemanticGroupNames.constants';
import {TSpan} from '../../types/Span.types';
import GenericSpanDetail from './components/GenericSpanDetail';

export interface ISpanDetailProps {
  testId?: string;
  span?: TSpan;
  resultId?: string;
}

const ComponentMap: Record<string, typeof GenericSpanDetail> = {
  [SemanticGroupNames.Http]: GenericSpanDetail,
};

const SpanDetail: React.FC<ISpanDetailProps> = ({span, ...props}) => {
  const Component = ComponentMap[span?.type || ''] || GenericSpanDetail;

  return <Component span={span} {...props} />;
};

export default SpanDetail;
