import {DownOutlined} from '@ant-design/icons';
import {Dropdown, Menu, Typography} from 'antd';
import {camelCase} from 'lodash';
import React from 'react';
import {DemoTestExampleList, IDemoTestExample} from 'constants/Test.constants';
import CreateTestAnalyticsService from 'services/Analytics/CreateTestAnalytics.service';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import * as S from './BasicDetails.styled';

interface IProps {
  onSelectDemo(demo: IDemoTestExample): void;
  selectedDemo?: IDemoTestExample;
}

const BasicDetailsDemoHelper = ({selectedDemo, onSelectDemo}: IProps) => {
  const handleOnDemoClick = ({key}: {key: string}) => {
    CreateTestAnalyticsService.onDemoTestClick();
    const demo = DemoTestExampleList.find(({name}) => name === key)!;
    onSelectDemo(demo);
  };

  return (
    <S.DemoContainer>
      <Typography.Text>Try these examples in our demo env: </Typography.Text>
      <Dropdown
        overlay={() => (
          <Menu
            items={DemoTestExampleList.map(({name}) => ({
              key: name,
              label: <span data-cy={`demo-example-${camelCase(name)}`}>{name}</span>,
            }))}
            onClick={handleOnDemoClick}
            style={{width: '190px'}}
          />
        )}
        placement="bottomLeft"
        trigger={['click']}
      >
        <Typography.Link data-cy="example-button" onClick={e => e.preventDefault()} strong>
          {selectedDemo?.name || 'Choose Example'} <DownOutlined />
        </Typography.Link>
      </Dropdown>
      <TooltipQuestion
        margin={8}
        title="We have a microservice based on the Pokemon API installed - pick an example to autofill this form"
      />
    </S.DemoContainer>
  );
};

export default BasicDetailsDemoHelper;
