import {Col, Form, Input, Row} from 'antd';
import {Body} from 'components/Inputs';
import {IPostmanValues, TDraftTestForm} from 'types/Test.types';
import {Auth, Headers, URL} from 'components/Fields';
import {CollectionFileField} from './fields/CollectionFileField';
import {EnvFileField} from './fields/EnvFileField';
import {SelectTestFromCollection} from './fields/SelectTestFromCollection';
import * as S from './UploadCollection.styled';

interface IProps {
  form: TDraftTestForm<IPostmanValues>;
}

const UploadCollectionForm = ({form}: IProps) => (
  <S.FieldsContainer>
    <Row gutter={12}>
      <Col span={18}>
        <Form.Item name="requests" hidden>
          <Input type="hidden" />
        </Form.Item>
        <Form.Item name="variables" hidden>
          <Input type="hidden" />
        </Form.Item>
        <CollectionFileField form={form} />
      </Col>
    </Row>
    <Row gutter={12}>
      <Col span={18}>
        <EnvFileField form={form} />
      </Col>
    </Row>
    <Row gutter={12}>
      <Col span={18}>
        <SelectTestFromCollection form={form} />
      </Col>
    </Row>
    <Row gutter={12}>
      <Col span={18}>
        <URL />
      </Col>
    </Row>
    <Row gutter={12}>
      <Col span={18}>
        <Headers />
      </Col>
    </Row>
    <Row gutter={12}>
      <Col span={18}>
        <Auth />
      </Col>
    </Row>
    <Row gutter={12}>
      <Col span={18}>
        <Form.Item name="body">
          <Body />
        </Form.Item>
      </Col>
    </Row>
  </S.FieldsContainer>
);

export default UploadCollectionForm;
