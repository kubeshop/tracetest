import {useEffect} from 'react';

export function useListenToArrowKeysDown(
  affectedSpans: string[],
  handleNextSpan: () => void,
  handlePrevSpan: () => void
): void {
  useEffect(() => {
    function downHandler({key}: KeyboardEvent) {
      if (affectedSpans.length > 0) {
        const focusedElement = document.activeElement;
        const advancedEditorId = 'advanced-editor';
        const isAdvancedEditor = focusedElement?.parentElement?.parentElement?.parentElement?.id === advancedEditorId;
        if (isAdvancedEditor) return;
        const searchElement = focusedElement?.tagName.toLowerCase() || '';
        if (!['input', 'textarea'].includes(searchElement)) {
          if (key === 'ArrowRight') {
            handleNextSpan();
          }
          if (key === 'ArrowLeft') {
            handlePrevSpan();
          }
        }
      }
    }

    window.addEventListener('keydown', downHandler);

    return () => {
      window.removeEventListener('keydown', downHandler);
    };
  }, [affectedSpans.length, handleNextSpan, handlePrevSpan]);
}
