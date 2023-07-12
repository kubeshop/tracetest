import {Space, Switch, Typography} from 'antd';
import {useState} from 'react';
import {LinterResultPlugin} from 'models/LinterResult.model';
import Trace from 'models/Trace.model';
import AnalyzerService from 'services/Analyzer.service';
import * as S from './AnalyzerResult.styled';
import AnalyzerScore from '../AnalyzerScore/AnalyzerScore';
import Rule from './Rule';
import Collapse, { CollapsePanel } from '../Collapse';

interface IProps {
  plugins: LinterResultPlugin[];
  trace: Trace;
}

const Plugins = ({plugins: rawPlugins, trace}: IProps) => {
  const [onlyErrors, setOnlyErrors] = useState(false);
  const plugins = AnalyzerService.getPlugins(rawPlugins, onlyErrors);

  return (
    <>
      <S.SwitchContainer>
        <Switch checked={onlyErrors} id="only_errors_enabled" onChange={() => setOnlyErrors(prev => !prev)} />
        <label htmlFor="only_errors_enabled">Show only errors</label>
      </S.SwitchContainer>

      <Collapse>
        {plugins.map(plugin => (
          <CollapsePanel
            header={
              <Space>
                <AnalyzerScore width="35px" height="35px" score={plugin.score} />
                <Typography.Text strong>{plugin.name}</Typography.Text>
                <Typography.Text type="secondary">{plugin.description}</Typography.Text>
              </Space>
            }
            key={plugin.name}
          >
            {plugin.rules.map(rule => (
              <Rule rule={rule} key={rule.name} trace={trace} />
            ))}
          </CollapsePanel>
        ))}
      </Collapse>
    </>
  );
};

export default Plugins;
