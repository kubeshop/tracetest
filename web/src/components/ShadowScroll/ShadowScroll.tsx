import {useCallback, useEffect, useLayoutEffect} from 'react';
import useShadowScroll from './hooks/useShadowScroll';
import * as S from './ShadowScroll.styled';

interface IProps {
  children: React.ReactNode;
}

const ShadowScroll = ({children}: IProps) => {
  const {
    contentRef,
    wrapperRef,
    topShadowRef,
    leftShadowRef,
    bottomShadowRef,
    rightShadowRef,
    updateHorizontalShadow,
    updateVerticalShadow,
    getIsReady,
  } = useShadowScroll();

  const onScroll = useCallback(
    ({target}: Event) => {
      const content = target as HTMLDivElement;
      updateHorizontalShadow(content);
      updateVerticalShadow(content);
    },
    [updateHorizontalShadow, updateVerticalShadow]
  );

  useEffect(() => {
    if (getIsReady()) {
      contentRef.current!.addEventListener('scroll', onScroll);
    }
  }, [contentRef, getIsReady, onScroll, wrapperRef]);

  useLayoutEffect(() => {
    if (getIsReady()) {
      updateHorizontalShadow(contentRef.current!);
      updateVerticalShadow(contentRef.current!);
    }
  }, [contentRef, getIsReady, updateHorizontalShadow, updateVerticalShadow]);

  return (
    <div style={{position: 'relative'}} ref={wrapperRef}>
      <S.Shadow $direction="to left" ref={rightShadowRef} />
      <S.Shadow $direction="to right" ref={leftShadowRef} />
      <S.Shadow $direction="to top" ref={bottomShadowRef} />
      <S.Shadow $direction="to bottom" ref={topShadowRef} />
      <S.Content ref={contentRef}>{children}</S.Content>
    </div>
  );
};

export default ShadowScroll;
