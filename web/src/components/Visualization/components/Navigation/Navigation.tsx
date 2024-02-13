import {LeftOutlined, RightOutlined} from '@ant-design/icons';
import {useCallback, useEffect} from 'react';
import {SupportedEditors} from 'constants/Editor.constants';

import * as S from './Navigation.styled';

const editorList = Object.values(SupportedEditors) as string[];

interface IProps {
  matchedSpans: string[];
  onNavigateToSpan(spanId: string): void;
  selectedSpan: string;
}

const Navigation = ({matchedSpans, onNavigateToSpan, selectedSpan}: IProps) => {
  // TODO: save matched spans in a different data structure
  const index = matchedSpans.findIndex(spanId => spanId === selectedSpan) + 1;

  const navigate = useCallback(
    (direction: 'next' | 'prev') => {
      const spanId =
        direction === 'next'
          ? matchedSpans[matchedSpans.indexOf(selectedSpan) + 1] || matchedSpans[0]
          : matchedSpans[matchedSpans.indexOf(selectedSpan) - 1] || matchedSpans[matchedSpans.length - 1];

      onNavigateToSpan(spanId);
    },
    [matchedSpans, onNavigateToSpan, selectedSpan]
  );

  useEffect(() => {
    if (!matchedSpans.length) return;

    function onKeydown({key}: KeyboardEvent) {
      const activeElement = document.activeElement;
      const activeElementTagName = activeElement?.tagName.toLowerCase() || '';
      const isAdvancedEditor = editorList.includes(
        activeElement?.parentElement?.parentElement?.parentElement?.id || ''
      );

      if (isAdvancedEditor || ['input', 'textarea'].includes(activeElementTagName)) return;

      if (key === 'ArrowRight') {
        navigate('next');
      }
      if (key === 'ArrowLeft') {
        navigate('prev');
      }
    }

    window.addEventListener('keydown', onKeydown);

    return () => {
      window.removeEventListener('keydown', onKeydown);
    };
  }, [matchedSpans.length, navigate]);

  if (!matchedSpans.length) return <div />;

  return (
    <S.Container>
      <S.ToggleButton id="span-back" onClick={() => navigate('prev')} icon={<LeftOutlined />} type="text" />
      <S.ToggleButton id="span-forward" onClick={() => navigate('next')} icon={<RightOutlined />} type="text" />
      <S.NavigationText strong>
        {index} of {matchedSpans.length} total
      </S.NavigationText>
    </S.Container>
  );
};

export default Navigation;
