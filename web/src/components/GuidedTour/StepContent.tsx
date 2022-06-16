import {CloseOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import React, {Dispatch} from 'react';
import {Container, Divider, Header, Title, TitleContainer} from './StepContent.styled';

interface IProps {
  title: string;
  setIsOpen?: Dispatch<React.SetStateAction<Boolean>>;
}

export const StepContent: React.FC<IProps> = ({setIsOpen, title, children}) => (
  <>
    <Header>
      <TitleContainer>
        <Title>{title}</Title>
      </TitleContainer>
      <Button
        style={{width: 24, height: 24}}
        type="link"
        icon={<CloseOutlined />}
        onClick={() => setIsOpen?.(o => !o)}
      />
    </Header>
    <Divider />
    <Container>{children}</Container>
    <Divider />
  </>
);
