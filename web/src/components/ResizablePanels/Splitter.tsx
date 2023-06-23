import {DoubleLeftOutlined, DoubleRightOutlined} from '@ant-design/icons';
import {Button} from 'antd';

import * as S from './ResizablePanels.styled';

interface IProps {
  isOpen: boolean;
  onClick(): void;
  id?: string;
  key: string | number;
  className?: string;
  onMouseDown: (e: React.MouseEvent<HTMLElement, MouseEvent>) => void;
  onTouchStart: (e: React.TouchEvent<HTMLElement>) => void;
}

const Splitter = ({isOpen, onClick, id, key, className, onMouseDown, onTouchStart}: IProps) => {
  return (
    <S.SplitterContainer id={id} key={key} className={className} onMouseDown={onMouseDown} onTouchStart={onTouchStart}>
      <S.ButtonContainer>
        <Button
          data-cy="toggle-drawer"
          icon={isOpen ? <DoubleLeftOutlined /> : <DoubleRightOutlined />}
          onClick={event => {
            event.stopPropagation();
            onClick();
          }}
          onMouseDown={event => event.stopPropagation()}
          shape="circle"
          size="small"
          type="primary"
        />
      </S.ButtonContainer>
    </S.SplitterContainer>
  );
};

export default Splitter;
