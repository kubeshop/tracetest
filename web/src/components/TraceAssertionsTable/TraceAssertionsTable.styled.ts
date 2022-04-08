import {Badge} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  margin-bottom: 36px;
`;

export const Header = styled.div`
  padding: 8px 0;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  width: 100%;
`;

export const LabelBadge = styled(Badge)`
  > sup {
    background-color: #f0f0f0;
    color: black;
    margin-left: 6px;
  }
`;

export const SelectedLabelBadge = styled(LabelBadge)`
  > sup {
    background-color: rgb(0,161,253);
    color: white;
  }
`;
