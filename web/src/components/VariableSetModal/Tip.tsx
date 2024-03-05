import {Typography} from 'antd';
import {VARIABLE_SET_DOCUMENTATION_URL} from 'constants/Common.constants';
import * as S from './VariableSetModal.styled';

const Tip = () => (
  <S.TipContainer>
    <Typography.Title level={3}>💡 What are Variable Sets?</Typography.Title>
    <Typography.Paragraph>
      <Typography.Text type="secondary">
        {`Variable sets are groups of variables that can be referenced by tests using this syntax: http://\${var:HOSTNAME}/api/users. `}
      </Typography.Text>
      <Typography.Link href={VARIABLE_SET_DOCUMENTATION_URL} target="_blank">
        Read more
      </Typography.Link>
    </Typography.Paragraph>
  </S.TipContainer>
);

export default Tip;
