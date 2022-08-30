import {DoubleLeftOutlined, DoubleRightOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import React, {useState} from 'react';

import * as S from './Drawer.styled';

interface IProps {
  children?: React.ReactNode;
}

const Drawer = ({children}: IProps) => {
  const [isAsideOpen, setIsAsideOpen] = useState(false);

  return (
    <S.Container $isOpen={isAsideOpen}>
      <S.Content>{children}</S.Content>
      <S.ButtonContainer>
        <Button
          icon={isAsideOpen ? <DoubleLeftOutlined /> : <DoubleRightOutlined />}
          onClick={() => setIsAsideOpen(isOpen => !isOpen)}
          shape="circle"
          size="small"
          type="primary"
        />
      </S.ButtonContainer>
    </S.Container>
  );
};

export default Drawer;
