import {ANALYZER_RULES_DOCUMENTATION_URL} from 'constants/Common.constants';
import TraceAnalyzerAnalytics from 'services/Analytics/TraceAnalyzer.service';
import * as S from './AnalyzerResult.styled';

interface IProps {
  id: string;
  isSmall?: boolean;
}

const RuleLink = ({id, isSmall = false}: IProps) => (
  <div>
    <S.RuleLinkText $isSmall={isSmall}>
      For more information, see{' '}
      <a
        href={`${ANALYZER_RULES_DOCUMENTATION_URL}/${id}`}
        onClick={() => TraceAnalyzerAnalytics.onDocsClick()}
        target="_blank"
      >
        analyzer({id})
      </a>
    </S.RuleLinkText>
  </div>
);

export default RuleLink;
