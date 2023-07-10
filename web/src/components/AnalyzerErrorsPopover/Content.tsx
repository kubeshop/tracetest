import RuleLink from 'components/AnalyzerResult/RuleLink';
import {TAnalyzerError} from 'types/TestRun.types';
import * as S from './AnalyzerErrorsPopover.styled';

interface IProps {
  errors: TAnalyzerError[];
}

const Content = ({errors}: IProps) => (
  <S.ContentContainer>
    {errors.map(analyzerError => (
      <S.RuleContainer key={analyzerError.ruleName}>
        <S.Text>{analyzerError.ruleName}</S.Text>

        {analyzerError.errors.length > 1 && (
          <>
            <div>
              <S.Text type="secondary">{analyzerError.ruleErrorDescription}</S.Text>
            </div>
            <S.List>
              {analyzerError.errors.map(error => (
                <li key={error.value}>
                  <S.Text type="secondary">{error.value}</S.Text>
                </li>
              ))}
            </S.List>
          </>
        )}

        {analyzerError.errors.length === 1 && (
          <div>
            <S.Text type="secondary">{analyzerError.errors[0].description}</S.Text>
          </div>
        )}

        {analyzerError.ruleDocumentationUrl && (
          <RuleLink id={analyzerError.ruleId} url={analyzerError.ruleDocumentationUrl} />
        )}
      </S.RuleContainer>
    ))}
  </S.ContentContainer>
);

export default Content;
