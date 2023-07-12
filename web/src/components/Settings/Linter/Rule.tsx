import {Col, Form, Input, Row, Select} from 'antd';
import {TDraftLinter} from 'types/Settings.types';
import {LinterRule, LinterRuleErrorLevel} from 'models/Linter.model';

interface IProps {
  fieldKey: number;
  baseName: string[];
  isDisabled: boolean;
}

const Rule = ({fieldKey, baseName, isDisabled}: IProps) => {
  const form = Form.useFormInstance<TDraftLinter>();
  const {name}: LinterRule = Form.useWatch([...baseName], form) ?? LinterRule({});

  return (
    <Row gutter={12} align="middle">
      <Col span={11}>
        <Form.Item name={[fieldKey, 'errorLevel']} label={name}>
          <Select disabled={isDisabled}>
            {Object.values(LinterRuleErrorLevel).map(level => {
              return (
                <Select.Option key={level} value={level}>
                  {level}
                </Select.Option>
              );
            })}
          </Select>
        </Form.Item>
      </Col>
      <Col span="auto">-</Col>
      <Col span={11}>
        <Form.Item
          name={[fieldKey, 'weight']}
          label="Weight input configuration"
          normalize={value => parseInt(String(value ?? 0), 10)}
        >
          <Input type="number" disabled={isDisabled} />
        </Form.Item>
      </Col>
    </Row>
  );
};

export default Rule;
