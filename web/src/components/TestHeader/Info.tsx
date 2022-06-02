import {Button, Popover, Typography} from 'antd';
import {InfoCircleOutlined} from '@ant-design/icons';

import Date from 'utils/Date';

interface IProps {
  date: string;
  executionTime: number;
  totalSpans: number;
  traceId: string;
}

const Info = ({date, executionTime, totalSpans, traceId}: IProps) => {
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
    <Popover placement="right" content={content}>
      <Button
        icon={<InfoCircleOutlined style={{color: 'rgba(3, 24, 73, 0.6)'}} />}
        shape="circle"
        size="small"
        type="text"
      />
    </Popover>
  );
};

export default Info;
