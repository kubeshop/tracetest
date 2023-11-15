import {DownOutlined} from '@ant-design/icons';
import {useCallback, useState} from 'react';
import {Dropdown, Form, Menu} from 'antd';
import {camelCase} from 'lodash';
import {IBasicValues, TDraftTest} from 'types/Test.types';
import CreateTestAnalyticsService from 'services/Analytics/CreateTestAnalytics.service';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import {useCreateTest} from 'providers/CreateTest';
import * as S from './DemoSelector.styled';

const BasicDetailsDemoHelper = () => {
  const [selectedDemo, setSelectedDemo] = useState<TDraftTest>();
  const form = Form.useFormInstance<IBasicValues>();

  const {demoList} = useCreateTest();

  const onSelectDemo = useCallback(
    ({key}: {key: string}) => {
      CreateTestAnalyticsService.onDemoTestClick();
      const demo = demoList.find(({name}) => name === key)!;

      form.setFieldsValue(demo);

      setSelectedDemo(demo);
    },
    [demoList, form]
  );

  return (
    <S.DemoContainer>
      <Dropdown
        overlay={() => (
          <Menu
            items={demoList.map(({name = ''}) => ({
              key: name,
              label: <span data-cy={`demo-example-${camelCase(name)}`}>{name}</span>,
            }))}
            onClick={onSelectDemo}
            style={{width: '190px'}}
          />
        )}
        placement="bottomLeft"
        trigger={['click']}
      >
        <S.Button data-cy="example-button" onClick={e => e.preventDefault()}>
          {selectedDemo?.name || 'Choose Example'} <DownOutlined />
        </S.Button>
      </Dropdown>
      <TooltipQuestion
        margin={8}
        title="We have a microservice based on the Pokemon API installed - pick an example to autofill this form"
      />
    </S.DemoContainer>
  );
};

export default BasicDetailsDemoHelper;
