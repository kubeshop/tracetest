import {Drawer} from 'antd';
import {noop} from 'lodash';
import {createContext, ReactNode, useCallback, useContext, useEffect, useMemo, useState} from 'react';
import {ReflexContainer, ReflexElement} from 'react-reflex';

import * as S from './RunLayout.styled';

enum OPEN_BOTTOM_PANEL_STATE {
  FORM,
  NORMAL,
}

const BOTTOM_PANEL_SIZES = {
  CLOSE: 85,
  INITIAL: 85,
  FORM: Math.round(window.innerWidth * 0.5),
  OPEN: Math.round(window.innerWidth * 0.5),
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
  //
  const openBottomPanel = useCallback((state: OPEN_BOTTOM_PANEL_STATE = OPEN_BOTTOM_PANEL_STATE.NORMAL) => {
    // setVisible(true);
  }, []);
  //
  const toggleBottomPanel = useCallback(() => {
    if (sizeBottomPanel <= BOTTOM_PANEL_SIZES.CLOSE) {
      setSizeBottomPanel(lastSizeBottomPanel || BOTTOM_PANEL_SIZES.OPEN);
      return;
    }
    setSizeBottomPanel(BOTTOM_PANEL_SIZES.CLOSE);
  }, [lastSizeBottomPanel, sizeBottomPanel]);
  //

  const value = useMemo(
    () => ({
      isBottomPanelOpen: sizeBottomPanel > BOTTOM_PANEL_SIZES.CLOSE,
      openBottomPanel,
      toggleBottomPanel,
    }),
    [openBottomPanel, sizeBottomPanel, toggleBottomPanel]
  );

  const [isResizing, setIsResizing] = useState(false);

  const onMouseDown = () => {
    setIsResizing(true);
  };

  const onMouseUp = () => {
    setIsResizing(false);
  };
  const onMouseMove = (e: any) => {
    if (isResizing) {
      let offsetRight = window.innerWidth - (e.clientX - document.body.offsetLeft);
      const minWidth = BOTTOM_PANEL_SIZES.CLOSE;
      const maxWidth = BOTTOM_PANEL_SIZES.OPEN;
      if (offsetRight > minWidth && offsetRight < maxWidth) {
        setSizeBottomPanel(offsetRight);
      }
    }
  };
  useEffect(() => {
    document.addEventListener('mousemove', onMouseMove);
    document.addEventListener('mouseup', onMouseUp);

    return () => {
      document.removeEventListener('mousemove', onMouseMove);
      document.removeEventListener('mouseup', onMouseUp);
    };
  });

  const style: React.CSSProperties = {
    position: 'absolute',
    width: '5px',
    padding: '4px 0 0',
    top: 0,
    left: 0,
    bottom: 0,
    zIndex: 100,
    cursor: 'ew-resize',
    backgroundColor: '#f4f7f9',
  };
  return (
    <RunLayoutContext.Provider value={value}>
      <div style={{display: 'flex', height: '100%'}}>
        <ReflexContainer orientation="vertical">
          <ReflexElement>
            <S.PanelContainer>{topPanel}</S.PanelContainer>
          </ReflexElement>
        </ReflexContainer>
        <Drawer
          style={{marginTop: 85}}
          closable={false}
          mask={false}
          placement="right"
          width={sizeBottomPanel}
          visible
          bodyStyle={{padding: 0, overflow: 'hidden'}}
        >
          <div style={style} onMouseDown={onMouseDown} />
          <div style={{display: 'flex', height: '100%'}}>{bottomPanel}</div>
        </Drawer>
      </div>
    </RunLayoutContext.Provider>
  );
};

export {OPEN_BOTTOM_PANEL_STATE, RunLayoutProvider, useRunLayout};
