import {FixedSizeList as List} from 'react-window';
import AutoSizer, {Size} from 'react-virtualized-auto-sizer';
import {Space, Tooltip, Typography} from 'antd';
import {PercentageOutlined} from '@ant-design/icons';
import {LinterResultPluginRule} from 'models/LinterResult.model';
import {LinterRuleErrorLevel} from 'models/Linter.model';
import * as S from './AnalyzerResult.styled';
import RuleIcon from './RuleIcon';
import RuleResult from './RuleResult';
import Collapse, {CollapsePanel} from '../Collapse';

interface IProps {
  rule: LinterResultPluginRule;
}

const Rule = ({rule: {tips, id, passed, description, name, level, results, weight = 0}, rule}: IProps) => {
  return (
    <Collapse>
      <CollapsePanel
        header={
          <S.Column>
            <S.RuleHeader>
              <Space>
                <RuleIcon passed={passed} level={level} />
                <Tooltip title={tips.join(' - ')}>
                  <Typography.Text strong>{name}</Typography.Text>
                </Tooltip>
                <Typography.Text type="secondary">{description}</Typography.Text>
              </Space>
            </S.RuleHeader>
            {level === LinterRuleErrorLevel.ERROR && (
              <Typography.Text type="secondary" style={{paddingLeft: 20}}>
                {weight}
                <PercentageOutlined />
              </Typography.Text>
            )}
          </S.Column>
        }
        key={id}
      >
        <S.RuleBody $resultCount={results.length}>
          <AutoSizer>
            {({height, width}: Size) => (
              <List height={height} itemCount={results.length} itemData={rule} itemSize={32} width={width}>
                {RuleResult}
              </List>
            )}
          </AutoSizer>
        </S.RuleBody>
      </CollapsePanel>
    </Collapse>
  );
};

export default Rule;
