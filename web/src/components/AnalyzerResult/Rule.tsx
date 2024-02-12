import {FixedSizeList as List} from 'react-window';
import AutoSizer, {Size} from 'react-virtualized-auto-sizer';
import {Space, Tooltip, Typography} from 'antd';
import {LinterResultPluginRule} from 'models/LinterResult.model';
import {LinterRuleErrorLevel} from 'models/Linter.model';
import * as S from './AnalyzerResult.styled';
import RuleIcon from './RuleIcon';
import RuleResult from './RuleResult';

interface IProps {
  rule: LinterResultPluginRule;
}

const Rule = ({rule: {tips, passed, description, name, level, results, weight = 0}, rule}: IProps) => {
  return (
    <S.RuleContainer>
      <S.Column>
        <S.RuleHeader>
          <Space>
            <RuleIcon passed={passed} level={level} />
            <Tooltip title={tips.join(' - ')}>
              <Typography.Text strong>{name}</Typography.Text>
            </Tooltip>
          </Space>
        </S.RuleHeader>
        <Typography.Text type="secondary" style={{paddingLeft: 20}}>
          {description}
        </Typography.Text>
        {level === LinterRuleErrorLevel.ERROR && (
          <Typography.Text type="secondary" style={{paddingLeft: 20}}>
            Weight: {weight}
          </Typography.Text>
        )}
      </S.Column>

      <S.RuleBody $resultCount={results.length}>
        <AutoSizer>
          {({height, width}: Size) => (
            <List height={height} itemCount={results.length} itemData={rule} itemSize={32} width={width}>
              {RuleResult}
            </List>
          )}
        </AutoSizer>
      </S.RuleBody>
    </S.RuleContainer>
  );
};

export default Rule;
