import {SelectOutlined} from '@ant-design/icons';
import {Button, Tooltip} from 'antd';

interface IProps {
  isValid: boolean;
  onSelect(): void;
}

const SpanActions = ({isValid, onSelect}: IProps) => (
  <Tooltip title="Select span">
    <Button disabled={!isValid} icon={<SelectOutlined />} onClick={() => onSelect()} size="small" type="link" />
  </Tooltip>
);

export default SpanActions;
