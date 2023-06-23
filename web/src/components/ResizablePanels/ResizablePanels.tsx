import * as Spaces from 'react-spaces';
import {useCallback, useMemo} from 'react';
import * as S from './ResizablePanels.styled';
import useResizablePanels, {TSize} from './hooks/useResizablePanels';
import Splitter from './Splitter';

export type TPanelComponentProps = {size: TSize};

export type TPanel = {
  name: string;
  isDefaultOpen?: boolean;
  minSize?: number;
  maxSize?: number;
  splitterPosition?: 'before' | 'after' | undefined;
  component(props: TPanelComponentProps): React.ReactElement;
};

interface IProps {
  panels: TPanel[];
}

const ResizablePanels = ({panels}: IProps) => {
  const {onStopResize, toggle, sizes} = useResizablePanels({panels});

  const getPanel = useCallback((name: string) => panels.find(panel => panel.name === name), [panels]);

  const elements = useMemo(
    () =>
      Object.values(sizes).reduce<React.ReactNode[]>((acc, size) => {
        const {component: Component, splitterPosition} = getPanel(size.name)!;

        if (splitterPosition === 'before') {
          acc.push(
            <Spaces.RightResizable
              onResizeEnd={newSize => onStopResize(size, newSize)}
              minimumSize={size.minSize}
              maximumSize={size.maxSize}
              size={size.size}
              key={size.name}
              handleRender={props => <Splitter {...props} isOpen={!size.isOpen} onClick={() => toggle(size)} />}
            >
              <Component size={size} />
            </Spaces.RightResizable>
          );
        } else if (splitterPosition === 'after') {
          acc.push(
            <Spaces.LeftResizable
              onResizeEnd={newSize => onStopResize(size, newSize)}
              minimumSize={size.minSize}
              maximumSize={size.maxSize}
              size={size.size}
              key={size.name}
              handleRender={props => <Splitter {...props} isOpen={size.isOpen} onClick={() => toggle(size)} />}
            >
              <Component size={size} />
            </Spaces.LeftResizable>
          );
        } else {
          acc.push(
            <Spaces.Fill>
              <Component size={size} />
            </Spaces.Fill>
          );
        }

        return acc;
      }, []),
    [getPanel, onStopResize, sizes, toggle]
  );

  return (
    <>
      <S.GlobalStyle />
      <Spaces.Fixed height="100%" width="100vw">
        {elements}
      </Spaces.Fixed>
    </>
  );
};

export default ResizablePanels;
