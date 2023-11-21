import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  background: ${({theme}) => theme.color.white};
  display: flex;
  flex-direction: column;
  margin-bottom: 24px;
  padding: 16px;
`;

export const FormContainer = styled.div`
  margin: 16px 0;
`;

export const FooterContainer = styled.div`
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
  padding: 16px 22px 0;
  border-top: 1px solid ${({theme}) => theme.color.borderLight};
`;

export const Description = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
  }
`;

export const LinterPluginsContainer = styled.div`
  margin-top: 10px;
  padding: 24px 24px 0;
  border-top: 1px solid ${({theme}) => theme.color.borderLight};
`;

export const LinterPluginList = styled.ul`
  && {
    list-style-type: disc;
  }
`;

export const SwitchContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
`;

export const SwitchLabel = styled.label<{$disabled?: boolean}>`
  color: ${({$disabled, theme}) => ($disabled ? theme.color.textLight : theme.color.text)};
  cursor: ${({$disabled}) => ($disabled ? 'not-allowed' : 'pointer')};
  margin-bottom: 24px;
`;
