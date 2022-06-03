import {QuestionCircleOutlined} from '@ant-design/icons';
import {Tooltip} from 'antd';

export const TooltipQuestion: React.FC<{title: string; margin?: number}> = ({margin = 8, title}) => (
  <Tooltip color="#FBFBFF" title={title}>
    <QuestionCircleOutlined style={{color: '#8C8C8C', cursor: 'pointer', marginLeft: margin}} />
  </Tooltip>
);
