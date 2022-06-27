import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  position: relative;
  height: 100%;

  > .react-flow {
    height: calc(100% - 76px);
  }
`;

export const SkeletonDiagramMessage = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  margin-top: 24px;
`;

export const SkeletonNode = styled.div`
  background-color: white;
  border: 1px solid #c9cedb;
  border-radius: 2px;
  min-width: fit-content;
  display: flex;
  width: 150px;
  max-width: 150px;
  height: 75px;
  justify-content: center;
  align-items: center;
`;

export const SkeletonNotch = styled.div`
  background-color: #fbfbff;
  position: absolute;
  top: 0;
  margin-top: 1px;
  padding: 3px 6px;
  border-radius: 2px 2px 0 0;
  width: 99%;
  font-weight: 700;
  height: 40px;
  padding-top: 10px;
`;

export const NameText = styled(Typography.Text).attrs({
  ellipsis: true,
})`
  margin: 0;
  font-size: 12px;
`;

export const TextContainer = styled.div`
  padding: 6px;
  padding-top: 30px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  justify-content: center;
  width: 150px;
  max-width: 150px;
  height: 38px;
  box-sizing: content-box;
`;

export const TextHolder = styled.div`
  @keyframes skeleton-loading {
    0% {
      background-color: hsl(200, 20%, 80%);
    }
    100% {
      background-color: hsl(200, 20%, 95%);
    }
  }

  animation: skeleton-loading 1s linear infinite alternate;
  height: 8px;
  border-radius: 2px;
  width: 100%;
`;
