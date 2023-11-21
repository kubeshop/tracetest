import styled from 'styled-components';

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
