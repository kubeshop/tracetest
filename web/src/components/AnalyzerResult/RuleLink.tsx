import * as S from './AnalyzerResult.styled';

interface IProps {
  id: string;
  url: string;
  isSmall?: boolean;
}

const RuleLink = ({id, url, isSmall = false}: IProps) => (
  <div>
    <S.RuleLinkText $isSmall={isSmall}>
      For more information, see{' '}
      <a href={url} target="_blank">
        analyzer({id})
      </a>
    </S.RuleLinkText>
  </div>
);

export default RuleLink;
