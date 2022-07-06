import styled from 'styled-components';

export const Container = styled.main`
  display: flex;
  height: 100%;
  width: 100%;
`;

export const LeftPanel = styled.div`
  display: flex;
  flex-basis: 50%;
  flex-direction: column;
  padding: 24px;
`;

export const RightPanel = styled.div`
  background: ${({theme}) => theme.color.white};
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
  flex-basis: 50%;
  overflow-y: scroll;
`;
