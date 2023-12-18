import {DownOutlined} from '@ant-design/icons';
import {useCallback, useState} from 'react';
import {Dropdown, Form, Menu} from 'antd';
import {camelCase} from 'lodash';
import {IBasicValues, TDraftTest} from 'types/Test.types';
import CreateTestAnalytics from 'services/Analytics/CreateTest.service';
import {TooltipQuestion} from 'components/TooltipQuestion/TooltipQuestion';
import * as S from './DemoSelector.styled';

interface IProps {
  demos: TDraftTest[];
}

const BasicDetailsDemoHelper = ({demos}: IProps) => {
  const [selectedDemo, setSelectedDemo] = useState<TDraftTest>();
  const form = Form.useFormInstance<IBasicValues>();

  const onSelectDemo = useCallback(
    ({key}: {key: string}) => {
      CreateTestAnalytics.onDemoClick();
      const demo = demos.find(({name}) => name === key)!;
      form.setFieldsValue(demo);
      setSelectedDemo(demo);
    },
    [demos, form]
  );

  return (
    <S.DemoContainer>
      <Dropdown
        overlay={() => (
          <Menu
            items={demos.map(({name = ''}) => ({
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
