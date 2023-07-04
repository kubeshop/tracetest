import {QuestionCircleOutlined} from '@ant-design/icons';
import {Tooltip} from 'antd';
import {useTheme} from 'styled-components';

interface IProps {
  margin?: number;
  title: React.ReactNode;
}

export const TooltipQuestion = ({margin = 8, title}: IProps) => {
  const theme = useTheme();

  return (
    <Tooltip title={title}>
      <QuestionCircleOutlined style={{color: theme.color.textSecondary, cursor: 'pointer', marginLeft: margin}} />
    </Tooltip>
  );
};
