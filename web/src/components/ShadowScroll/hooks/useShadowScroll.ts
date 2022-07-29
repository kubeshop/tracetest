import {useCallback, useRef} from 'react';

const useShadowScroll = () => {
  const wrapperRef = useRef<HTMLDivElement>(null);
  const contentRef = useRef<HTMLDivElement>(null);
  const rightShadowRef = useRef<HTMLDivElement>(null);
  const leftShadowRef = useRef<HTMLDivElement>(null);

  const bottomShadowRef = useRef<HTMLDivElement>(null);
  const topShadowRef = useRef<HTMLDivElement>(null);

  const getIsReady = useCallback(
    () => Boolean(wrapperRef.current && contentRef.current && rightShadowRef.current && leftShadowRef.current),
    []
  );

  const getScrollWidth = useCallback(
    () =>
      (contentRef.current && wrapperRef.current && contentRef.current.scrollWidth - wrapperRef.current.offsetWidth) ||
      1,
    [wrapperRef, contentRef]
  );

  const getScrollHeight = useCallback(
    () =>
      (contentRef.current && wrapperRef.current && contentRef.current.scrollHeight - wrapperRef.current.offsetHeight) ||
      1,
    [wrapperRef, contentRef]
  );

  const updateHorizontalShadow = useCallback(
    ({scrollLeft}: HTMLDivElement) => {
      const scrollWidth = getScrollWidth();
      const currentScroll = scrollLeft / scrollWidth;

      rightShadowRef!.current!.style.opacity = `${1 - currentScroll}`;
      leftShadowRef!.current!.style.opacity = `${currentScroll}`;
    },
    [getScrollWidth]
  );

  const updateVerticalShadow = useCallback(
    ({scrollTop}: HTMLDivElement) => {
      const scrollHeight = getScrollHeight();
      const currentScroll = scrollTop / scrollHeight;

      bottomShadowRef!.current!.style.opacity = `${1 - currentScroll}`;
      topShadowRef!.current!.style.opacity = `${currentScroll}`;
    },
    [getScrollHeight]
  );

  return {
    updateHorizontalShadow,
    updateVerticalShadow,
    getIsReady,
    contentRef,
    wrapperRef,
    rightShadowRef,
    leftShadowRef,
    bottomShadowRef,
    topShadowRef,
  };
};

export default useShadowScroll;
