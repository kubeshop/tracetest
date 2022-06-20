import {DownOutlined} from '@ant-design/icons';
import {Dropdown, FormInstance, Menu, Typography} from 'antd';
import {camelCase} from 'lodash';
import React from 'react';
import {DemoTestExampleList} from '../../constants/Test.constants';
import CreateTestAnalyticsService from '../../services/Analytics/CreateTestAnalytics.service';
import {TooltipQuestion} from '../TooltipQuestion/TooltipQuestion';
import {ICreateTestValues} from './CreateTestForm';
import * as S from './CreateTestModal.styled';

interface IProps {
  form: FormInstance<ICreateTestValues>;
  onSelectDemo(value: string): void;
  onValidation(isValid: boolean): void;
  selectedDemo: string;
}

export const CreateTestFormDemoHelper: React.FC<IProps> = ({selectedDemo, onSelectDemo, onValidation, form}) => {
  const handleOnDemoClick = ({key}: {key: string}) => {
    CreateTestAnalyticsService.onDemoTestClick();
    onSelectDemo(key);
    const {body, description, method, name, url} = DemoTestExampleList.find(demo => demo.name === key) || {};

    form.setFieldsValue({
      body,
      method,
      name: `${name} - ${description}`,
      url,
    });

    onValidation(true);
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
          {selectedDemo || 'Choose Example'} <DownOutlined />
        </Typography.Link>
      </Dropdown>
      <TooltipQuestion
        margin={8}
        title="We have a microservice based on the Pokemon API installed - pick an example to autofill this form"
      />
    </S.DemoContainer>
  );
};
