import styled, {css} from 'styled-components';

export const DemoContainer = styled.div`
  margin-bottom: 24px;
`;

export const InputContainer = styled.div<{$isEditing?: boolean}>`
  display: grid;
  gap: 26px;
  grid-template-columns: 55%;

  ${({$isEditing = false}) =>
    $isEditing &&
    css`
      grid-template-columns: 100%;
    `};
`;
