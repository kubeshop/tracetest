import {DoubleLeftOutlined, DoubleRightOutlined} from '@ant-design/icons';
import {Tooltip, TooltipProps} from 'antd';
import * as S from './ResizablePanels.styled';

interface IProps {
  isOpen: boolean;
  onClick(): void;
  id?: string;
  className?: string;
  name: string;
  isToolTipVisible?: boolean;
  tooltip?: string;
  tooltipPlacement?: TooltipProps['placement'];
  onMouseDown(e: React.MouseEvent<HTMLElement, MouseEvent>): void;
  onTouchStart(e: React.TouchEvent<HTMLElement>): void;
  dataTour?: string;
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
  isToolTipVisible = false,
  dataTour,
}: IProps) => {
  const button = (
    <S.SplitterButton
      $isPulsing={isToolTipVisible}
      data-cy={`toggle-drawer-${name}`}
      icon={isOpen ? <DoubleLeftOutlined /> : <DoubleRightOutlined />}
      onClick={event => {
        event.stopPropagation();
        onClick();
      }}
      onMouseDown={event => event.stopPropagation()}
      shape="circle"
      type="primary"
      data-tour={dataTour}
    />
  );

  return (
    <S.SplitterContainer id={id} key={id} className={className} onMouseDown={onMouseDown} onTouchStart={onTouchStart}>
      <S.ButtonContainer>
        {isToolTipVisible ? (
          <Tooltip title={tooltip} visible trigger={[]} placement={tooltipPlacement} overlayClassName="splitter">
            {button}
          </Tooltip>
        ) : (
          button
        )}
      </S.ButtonContainer>
    </S.SplitterContainer>
  );
};

export default Splitter;
