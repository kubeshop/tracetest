import {DoubleLeftOutlined, DoubleRightOutlined} from '@ant-design/icons';
import {Tooltip, TooltipProps} from 'antd';
import * as S from './ResizablePanels.styled';

interface IProps {
  isOpen: boolean;
  onClick(): void;
  id?: string;
  className?: string;
  name: string;
  tooltip?: string;
  tooltipPlacement?: TooltipProps['placement'];
  onMouseDown(e: React.MouseEvent<HTMLElement, MouseEvent>): void;
  onTouchStart(e: React.TouchEvent<HTMLElement>): void;
}

const Splitter = ({
  isOpen,
  name,
  onClick,
  id,
  className,
  onMouseDown,
  onTouchStart,
  tooltip,
  tooltipPlacement = 'right',
}: IProps) => (
  <S.SplitterContainer id={id} key={id} className={className} onMouseDown={onMouseDown} onTouchStart={onTouchStart}>
    <S.ButtonContainer>
      <Tooltip title={tooltip} trigger="hover" placement={tooltipPlacement} overlayClassName="splitter">
        <S.SplitterButton
          data-cy={`toggle-drawer-${name}`}
          icon={isOpen ? <DoubleLeftOutlined /> : <DoubleRightOutlined />}
          onClick={event => {
            event.stopPropagation();
            onClick();
          }}
          onMouseDown={event => event.stopPropagation()}
          shape="circle"
          type="primary"
        />
      </Tooltip>
    </S.ButtonContainer>
  </S.SplitterContainer>
);

export default Splitter;
