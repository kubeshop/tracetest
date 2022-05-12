import styled from 'styled-components';

export const TestDetailsHeader = styled.div`
  display: flex;
  width: 100%;
  justify-content: space-between;
  margin: 32px 0px 24px;
`;

export const Wrapper = styled.div<{detail?: boolean}>`
  padding: 0px 24px;
  display: flex;
  flex-grow: 1;
  flex-direction: column;
  background: #fff;
`;
