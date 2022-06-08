import {noop} from 'lodash';
import {createContext, ReactNode, useCallback, useContext, useMemo, useState} from 'react';
import {HandlerProps, ReflexContainer, ReflexElement, ReflexSplitter} from 'react-reflex';

import * as S from './RunLayout.styled';

enum OPEN_BOTTOM_PANEL_STATE {
  FORM,
  NORMAL,
}

const BOTTOM_PANEL_MIN_SIZE = 64;
const TOP_PANEL_MIN_SIZE = 100;

const BOTTOM_PANEL_SIZES = {
  /* closed size */
  CLOSE: BOTTOM_PANEL_MIN_SIZE,
  /* add assertion form size */
  FORM: 400,
  /* initial size on page load */
  INITIAL: Math.round(window.innerHeight * 0.2),
  /* default open size */
  OPEN: Math.round(window.innerHeight * 0.5),
};

interface IContext {
  isBottomPanelOpen: boolean;
  openBottomPanel(state?: OPEN_BOTTOM_PANEL_STATE): void;
  toggleBottomPanel(): void;
}

const RunLayoutContext = createContext<IContext>({
  isBottomPanelOpen: false,
  openBottomPanel: noop,
  toggleBottomPanel: noop,
});

const useRunLayout = () => {
  const context = useContext(RunLayoutContext);
  if (context === undefined) {
    throw new Error(`useRunLayout must be used within a RunLayoutProvider`);
  }
  return context;
};

interface IProps {
  bottomPanel: ReactNode;
  topPanel: ReactNode;
}

const RunLayoutProvider = ({bottomPanel, topPanel}: IProps) => {
  const [sizeBottomPanel, setSizeBottomPanel] = useState(BOTTOM_PANEL_SIZES.INITIAL);
  const [lastSizeBottomPanel, setLastSizeBottomPanel] = useState(0);

  const openBottomPanel = useCallback(
    (state: OPEN_BOTTOM_PANEL_STATE = OPEN_BOTTOM_PANEL_STATE.NORMAL) => {
      if (state === OPEN_BOTTOM_PANEL_STATE.FORM && sizeBottomPanel < BOTTOM_PANEL_SIZES.FORM) {
        setSizeBottomPanel(BOTTOM_PANEL_SIZES.FORM);
      }
      if (state === OPEN_BOTTOM_PANEL_STATE.NORMAL) {
        setSizeBottomPanel(lastSizeBottomPanel || BOTTOM_PANEL_SIZES.OPEN);
      }
    },
    [lastSizeBottomPanel, sizeBottomPanel]
  );

  const toggleBottomPanel = useCallback(() => {
    if (sizeBottomPanel <= BOTTOM_PANEL_SIZES.CLOSE) {
      setSizeBottomPanel(lastSizeBottomPanel || BOTTOM_PANEL_SIZES.OPEN);
      return;
    }
    setSizeBottomPanel(BOTTOM_PANEL_SIZES.CLOSE);
  }, [lastSizeBottomPanel, sizeBottomPanel]);

  const handleOnStopResize = useCallback(({domElement}: HandlerProps) => {
    const element = domElement as HTMLElement;
    const size = Math.round(element?.offsetHeight ?? 0);
    setSizeBottomPanel(size);
    setLastSizeBottomPanel(size <= BOTTOM_PANEL_SIZES.CLOSE ? 0 : size);
  }, []);

  const value = useMemo(
    () => ({
      isBottomPanelOpen: sizeBottomPanel > BOTTOM_PANEL_SIZES.CLOSE,
      openBottomPanel,
      toggleBottomPanel,
    }),
    [openBottomPanel, sizeBottomPanel, toggleBottomPanel]
  );

  return (
    <RunLayoutContext.Provider value={value}>
      <ReflexContainer>
        <ReflexElement minSize={TOP_PANEL_MIN_SIZE}>
          <S.PanelContainer>{topPanel}</S.PanelContainer>
        </ReflexElement>

        <ReflexSplitter
          style={{
            borderBottom: '1px solid #eeeeee',
            borderTop: '1px solid #eeeeee',
            boxShadow: '0px -4px 14px rgba(153, 155, 168, 0.25)',
          }}
        />

        <ReflexElement
          direction={-1}
          minSize={BOTTOM_PANEL_MIN_SIZE}
          onStopResize={handleOnStopResize}
          size={sizeBottomPanel}
        >
          <S.PanelContainer>{bottomPanel}</S.PanelContainer>
        </ReflexElement>
      </ReflexContainer>
    </RunLayoutContext.Provider>
  );
};

export {OPEN_BOTTOM_PANEL_STATE, RunLayoutProvider, useRunLayout};
