import {PanelResizeHandle} from 'react-resizable-panels';
import {DoubleLeftOutlined, DoubleRightOutlined} from '@ant-design/icons';

import * as S from './ResizablePanels.styled';

interface IProps {
  isOpen: boolean;
  onToggle(): void;
  placement?: 'left' | 'right';
}

const Handle = ({onToggle, isOpen, placement = 'left'}: IProps) => {
  return (
    <PanelResizeHandle className="panel-handle">
      <S.ButtonContainer $placement={placement}>
        <S.ToggleButton
          onClick={e => {
            e.stopPropagation();
            onToggle();
          }}
          icon={isOpen ? <DoubleLeftOutlined /> : <DoubleRightOutlined />}
        />
      </S.ButtonContainer>
    </PanelResizeHandle>
  );
};

export default Handle;
