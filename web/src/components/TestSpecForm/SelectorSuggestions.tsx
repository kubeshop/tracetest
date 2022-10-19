import {InfoCircleOutlined} from '@ant-design/icons';
import {ISuggestion} from 'types/TestSpecs.types';
import * as S from './TestSpecForm.styled';

interface IProps {
  onClick(query: string): void;
  selectorSuggestions: ISuggestion[];
}

const SelectorSuggestions = ({onClick, selectorSuggestions}: IProps) => (
  <>
    {Boolean(selectorSuggestions.length) && (
      <S.FormSectionText>
        <InfoCircleOutlined /> Try to select spans using these suggestions:
      </S.FormSectionText>
    )}
    {selectorSuggestions.map(selectorSuggestion => (
      <S.SuggestionsButton key={selectorSuggestion.query} onClick={() => onClick(selectorSuggestion.query)} type="link">
        {selectorSuggestion.title}
      </S.SuggestionsButton>
    ))}
  </>
);

export default SelectorSuggestions;
