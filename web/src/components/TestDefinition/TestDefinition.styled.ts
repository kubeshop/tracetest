import {Button, Typography} from 'antd';
import styled from 'styled-components';

export const CodeContainer = styled.div`
  position: relative;
  border: ${({theme}) => `1px solid ${theme.color.border}`};
  min-height: 370px;

  pre {
    margin: 0;
    min-height: inherit;
    max-height: 340px;
    background: ${({theme}) => theme.color.background} !important;
  }
`;

export const Container = styled.div`
  padding: 24px;
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin-bottom: 16px;
    font-weight: 700;
  }
`;

export const SubtitleContainer = styled.div`
  margin-bottom: 8px;
`;

export const CopyButton = styled(Button)`
  && {
    background-color: ${({theme}) => theme.color.white};
    padding-left: 8px;
    padding-right: 8px;
    font-size: ${({theme}) => theme.size.md};
    font-weight: 600;
  }
`;

export const CopyIconContainer = styled.div`
  position: absolute;
  right: 8px;
  top: 9px;
  padding: 0 2px;
  border-radius: 2px;
  cursor: pointer;
  background-color: ${({theme}) => theme.color.textHighlight};
  z-index: 101;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
`;
