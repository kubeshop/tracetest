import {Drawer} from 'antd';
import * as React from 'react';
import {MouseEventHandler, useCallback, useEffect, useState} from 'react';
import styled from 'styled-components';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import {useReferenceMemo} from './useReferenceMemo';

interface IProps {
  visiblePortion: number;
  children: JSX.Element[];
}

export enum DrawerState {
  OPEN = 'OPEN',
  CLOSE = 'CLOSE',
  INITIAL = 'INITIAL',
  FORM = 'FORM',
}

const CustomDrawer = styled(Drawer)`
  .ant-drawer-body {
    display: flex;
    flex-direction: column;
  }
`;

const ResizableDrawer: React.FC<IProps> = ({children, visiblePortion}: IProps) => {
  const {drawerState} = useAssertionForm();
  const [isResizing, setIsResizing] = useState(false);
  const ref = useReferenceMemo(visiblePortion);
  const [height, setHeight] = useState(ref[DrawerState.INITIAL]);

  const onPointerDown: MouseEventHandler = useCallback(() => {
    setIsResizing(true);
    window.addEventListener('selectstart', e => e.preventDefault());
  }, [setIsResizing]);

  const onMouseUp: EventListener = useCallback(() => {
    setIsResizing(false);
    window.removeEventListener('selectstart', e => e.preventDefault());
  }, [setIsResizing]);

  const onMouseMove: EventListener = useCallback(
    (e: MouseEventInit) => {
      if (isResizing) {
        const offsetRight =
          document.body.offsetHeight - ((e.clientY || document.body.offsetLeft) - document.body.offsetLeft);
        if (offsetRight > visiblePortion && offsetRight < ref['OPEN']) {
          setHeight(offsetRight);
        }
      }
    },
    [setHeight, isResizing, visiblePortion, ref]
  );

  useEffect(() => {
    if (drawerState === DrawerState.CLOSE) {
      setHeight(visiblePortion);
      return;
    }
    if (drawerState === DrawerState.OPEN) {
      setHeight(ref[DrawerState.OPEN]);
    }
    if (drawerState === DrawerState.FORM) {
      setHeight(ref[DrawerState.FORM]);
    }
    if (drawerState === DrawerState.INITIAL) {
      setHeight(ref[DrawerState.INITIAL]);
    }
  }, [drawerState, ref, visiblePortion]);
  useEffect(() => {
    document.addEventListener('pointermove', onMouseMove);
    document.addEventListener('pointerup', onMouseUp);

    return () => {
      document.removeEventListener('pointermove', onMouseMove);
      document.removeEventListener('pointerup', onMouseUp);
    };
  });

  return (
    <CustomDrawer
      placement="bottom"
      closable={false}
      visible
      height={height}
      mask={false}
      style={{overflow: 'hidden'}}
      bodyStyle={{overflow: 'hidden', padding: 0}}
    >
      <div
        id="draggg"
        style={{
          position: 'absolute',
          width: '100%',
          height: 5,
          padding: '4px 0 0',
          top: 0,
          left: 0,
          bottom: 0,
          zIndex: 100,
          cursor: 'row-resize',
          backgroundColor: '#f4f7f9',
        }}
        onPointerDown={onPointerDown}
      />
      {children.map(child => React.cloneElement(child, {height, max: ref['OPEN'], min: visiblePortion}))}
    </CustomDrawer>
  );
};

export default ResizableDrawer;
