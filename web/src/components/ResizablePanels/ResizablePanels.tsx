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
  position?: 'left' | 'right' | undefined;
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
      Object.values(sizes).reduce<React.ReactNode[]>((acc, size, index) => {
        const {component: Component, position} = getPanel(size.name)!;

        if (position === 'left') {
          return acc.concat([
            <Spaces.LeftResizable
              onResizeEnd={newSize => onStopResize(size, newSize)}
              minimumSize={size.minSize}
              maximumSize={size.maxSize}
              size={size.size}
              key={size.name}
              handleRender={props => (
                <Splitter {...props} name={size.name} isOpen={size.isOpen} onClick={() => toggle(size)} />
              )}
              order={index + 1}
            >
              <Component size={size} />
            </Spaces.LeftResizable>,
          ]);
        }

        if (position === 'right') {
          return acc.concat([
            <Spaces.RightResizable
              onResizeEnd={newSize => onStopResize(size, newSize)}
              minimumSize={size.minSize}
              maximumSize={size.maxSize}
              size={size.size}
              key={size.name}
              handleRender={props => (
                <Splitter {...props} name={size.name} isOpen={!size.isOpen} onClick={() => toggle(size)} />
              )}
              order={index + 1}
            >
              <Component size={size} />
            </Spaces.RightResizable>,
          ]);
        }

        return acc.concat(
          <Spaces.Fill>
            <Component size={size} />
          </Spaces.Fill>
        );
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
