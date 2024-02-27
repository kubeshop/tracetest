import {ArrowsAltOutlined, CloseCircleOutlined, ShrinkOutlined} from '@ant-design/icons';
import * as S from './ResizablePanels.styled';
import {TSize} from './hooks/useResizablePanel';

interface IProps {
  size: TSize;
  onChange(width: number): void;
}

const PanelContainer: React.FC<IProps> = ({
  children,
  size: {isOpen, fullScreen, openSize, closeSize, size},
  onChange,
}) => {
  return (
    <S.PanelContainer $isOpen={isOpen} onClick={() => !isOpen && onChange(openSize())}>
      <S.Controls $isOpen={isOpen}>
        {size < fullScreen() ? (
          <ArrowsAltOutlined onClick={() => onChange(fullScreen())} />
        ) : (
          <ShrinkOutlined onClick={() => onChange(openSize())} />
        )}
        <CloseCircleOutlined onClick={() => onChange(closeSize())} />
      </S.Controls>
      {children}
    </S.PanelContainer>
  );
};

export default PanelContainer;
