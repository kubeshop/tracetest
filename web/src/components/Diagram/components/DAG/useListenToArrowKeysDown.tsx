import {useEffect} from 'react';

export function useListenToArrowKeysDown(
  affectedSpans: string[],
  handleNextSpan: () => void,
  handlePrevSpan: () => void
): void {
  useEffect(() => {
    function downHandler({key}: KeyboardEvent) {
      if (affectedSpans.length > 0) {
        const searchElement = document.activeElement?.tagName.toLowerCase() || '';
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
