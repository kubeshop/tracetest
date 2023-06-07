import {DoubleLeftOutlined, DoubleRightOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import {noop} from 'lodash';
import {createContext, ReactNode, useCallback, useContext, useMemo, useState} from 'react';
import {HandlerProps, ReflexContainer, ReflexElement, ReflexSplitter} from 'react-reflex';
import TestRunAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {StepsID} from 'components/GuidedTour/testRunSteps';

import * as S from './Drawer.styled';

const LEFT_PANEL_MIN_SIZE = 30;
const RIGHT_PANEL_MIN_SIZE = 400;

const LEFT_PANEL_SIZES = {
  CLOSE: LEFT_PANEL_MIN_SIZE,
  OPEN: 270,
};

interface IProps {
  leftPanel: ReactNode;
  rightPanel: ReactNode;
}

interface DrawerContext {
  openDrawer: () => void;
}

export const Context = createContext<DrawerContext>({
  openDrawer: noop,
});

export const useDrawer = () => useContext(Context);

const Drawer = ({leftPanel, rightPanel}: IProps) => {
  const [sizeLeftPanel, setSizeLeftPanel] = useState(LEFT_PANEL_SIZES.CLOSE);
  const [lastSizeLeftPanel, setLastSizeLeftPanel] = useState(0);

  const toggleLeftPanel = useCallback(
    lastSize => {
      if (sizeLeftPanel <= LEFT_PANEL_SIZES.CLOSE) {
        setSizeLeftPanel(lastSize || LEFT_PANEL_SIZES.OPEN);
        return;
      }
      setSizeLeftPanel(LEFT_PANEL_SIZES.CLOSE);
    },
    [sizeLeftPanel]
  );

  const handleOnStopResize = useCallback(({domElement}: HandlerProps) => {
    const element = domElement as HTMLElement;
    const size = Math.round(element?.offsetWidth ?? 0);
    setSizeLeftPanel(size);
    setLastSizeLeftPanel(size <= LEFT_PANEL_SIZES.CLOSE ? 0 : size);
  }, []);

  const isOpen = sizeLeftPanel > LEFT_PANEL_SIZES.CLOSE;

  return (
    <>
      <S.GlobalStyle />
      <ReflexContainer orientation="vertical">
        <ReflexElement minSize={LEFT_PANEL_MIN_SIZE} onStopResize={handleOnStopResize} size={sizeLeftPanel}>
          <S.Content $isOpen={isOpen}>{leftPanel}</S.Content>
        </ReflexElement>

        <ReflexSplitter>
          <S.ButtonContainer>
            <Button
              data-cy="toggle-drawer"
              data-tour={StepsID.SpanDetails}
              icon={isOpen ? <DoubleLeftOutlined /> : <DoubleRightOutlined />}
              onClick={event => {
                event.stopPropagation();
                toggleLeftPanel(lastSizeLeftPanel);
                TestRunAnalyticsService.onAttributeDrawerOpen();
              }}
              onMouseDown={event => event.stopPropagation()}
              shape="circle"
              size="small"
              type="primary"
            />
          </S.ButtonContainer>
        </ReflexSplitter>

        <ReflexElement minSize={RIGHT_PANEL_MIN_SIZE}>
          <Context.Provider
            value={useMemo(
              () => ({
                openDrawer: () => {
                  if (!isOpen) toggleLeftPanel(lastSizeLeftPanel);
                },
              }),
              [isOpen, toggleLeftPanel, lastSizeLeftPanel]
            )}
          >
            {rightPanel}
          </Context.Provider>
        </ReflexElement>
      </ReflexContainer>
    </>
  );
};

export default Drawer;
