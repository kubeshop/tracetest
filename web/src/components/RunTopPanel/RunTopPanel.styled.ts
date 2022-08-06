import {ReflexContainer, ReflexElement} from 'react-reflex';
import styled from 'styled-components';

export const Container = styled(ReflexContainer)`
  display: flex;
  height: 100%;
  width: 100%;

  &.vertical > .reflex-splitter {
    border-color: transparent;
    background: transparent;
  }
`;

export const LeftPanel = styled(ReflexElement)`
  display: flex;
  flex-direction: column;
  //padding: 24px;
`;

export const RightPanel = styled(ReflexElement)`
  background: ${({theme}) => theme.color.white};
  box-shadow: 0 20px 24px rgba(153, 155, 168, 0.18);
  overflow-y: scroll;
`;
