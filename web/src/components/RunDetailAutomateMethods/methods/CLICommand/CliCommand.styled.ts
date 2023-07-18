import {Typography} from 'antd';
import styled from 'styled-components';

export const Title = styled(Typography.Title).attrs({
  level: 3,
})`
  && {
    font-size: ${({theme}) => theme.size.md};
    font-weight: 600;
    margin-bottom: 16px;
  }
`;

export const Subtitle = styled(Typography.Paragraph)`
  && {
    margin-bottom: 8px;
  }
`;

export const TitleContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const Container = styled.div`
  margin: 16px 0;
`;

export const SwitchContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
`;

export const SwitchLabel = styled.label<{$disabled?: boolean}>`
  color: ${({$disabled, theme}) => ($disabled ? theme.color.textLight : theme.color.text)};
  cursor: ${({$disabled}) => ($disabled ? 'not-allowed' : 'pointer')};
`;

export const ControlsContainer = styled.div`
  margin-top: 16px;
`;

export const OptionsContainer = styled.div`
  margin-bottom: 24px;
`;

export const FormatContainer = styled.div``;
