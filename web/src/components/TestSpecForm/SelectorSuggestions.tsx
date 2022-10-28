import {ArrowLeftOutlined, InfoCircleOutlined} from '@ant-design/icons';
import {Tooltip} from 'antd';

import {ISuggestion} from 'types/TestSpecs.types';
import * as S from './TestSpecForm.styled';

interface IProps {
  onClick(query: string): void;
  onClickPrevSelector(query: string): void;
  prevSelector: string;
  selectorSuggestions: ISuggestion[];
}

const SelectorSuggestions = ({onClick, onClickPrevSelector, prevSelector, selectorSuggestions}: IProps) => (
  <>
    {Boolean(selectorSuggestions.length) && (
      <>
        <S.FormSectionText>
          <InfoCircleOutlined /> Try to select spans using these suggestions:
        </S.FormSectionText>

        {Boolean(prevSelector) && (
          <Tooltip title={prevSelector}>
            <S.SuggestionsButton onClick={() => onClickPrevSelector(prevSelector)} type="link">
              <ArrowLeftOutlined /> Prev selector
            </S.SuggestionsButton>
          </Tooltip>
        )}
      </>
    )}

    {selectorSuggestions.map(selectorSuggestion => (
      <Tooltip key={selectorSuggestion.query} title={selectorSuggestion.query || '(empty selector)'}>
        <S.SuggestionsButton onClick={() => onClick(selectorSuggestion.query)} type="link">
          {selectorSuggestion.title}
        </S.SuggestionsButton>
      </Tooltip>
    ))}
  </>
);

export default SelectorSuggestions;
