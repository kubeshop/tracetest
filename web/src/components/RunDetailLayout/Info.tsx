import {Button, Popover, Typography} from 'antd';
import {InfoCircleOutlined} from '@ant-design/icons';
import {useTheme} from 'styled-components';

import Date from 'utils/Date';
import TestAnalyticsService from '../../services/Analytics/TestAnalytics.service';

interface IProps {
  date: string;
  executionTime: number;
  totalSpans: number;
  traceId: string;
}

const Info = ({date, executionTime, totalSpans, traceId}: IProps) => {
  const theme = useTheme();

  const content = (
    <>
      <div>
        <Typography.Text strong>Trace ID: </Typography.Text>
        <Typography.Text>{traceId}</Typography.Text>
      </div>
      <div>
        <Typography.Text strong>Trace transaction occurred: </Typography.Text>
        <Typography.Text>{Date.format(date, "yyyy/MM/dd 'at' HH:mm:ss")}</Typography.Text>
      </div>
      <div>
        <Typography.Text strong>Execution time: </Typography.Text>
        <Typography.Text>{executionTime} s</Typography.Text>
      </div>
      <div>
        <Typography.Text strong>Total spans: </Typography.Text>
        <Typography.Text>{totalSpans}</Typography.Text>
      </div>
    </>
  );

  return (
    <Popover
      placement="right"
      content={content}
      onVisibleChange={isVisible => {
        isVisible && TestAnalyticsService.onDisplayTestInfo();
      }}
    >
      <Button
        icon={<InfoCircleOutlined style={{color: theme.color.primary}} />}
        shape="circle"
        size="small"
        type="text"
      />
    </Popover>
  );
};

export default Info;
