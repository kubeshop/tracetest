import {Col, Form, Input, Row} from 'antd';
import {IPostmanValues, TDraftTestForm} from 'types/Test.types';
import RequestDetailsAuthInput from '../../../Rest/steps/RequestDetails/RequestDetailsAuthInput/RequestDetailsAuthInput';
import RequestDetailsHeadersInput from '../../../Rest/steps/RequestDetails/RequestDetailsHeadersInput';
import RequestDetailsUrlInput from '../../../Rest/steps/RequestDetails/RequestDetailsUrlInput';
import {CollectionFileField} from './fields/CollectionFileField';
import {EnvFileField} from './fields/EnvFileField';
import {SelectTestFromCollection} from './fields/SelectTestFromCollection';

export const FORM_ID = 'upload-collection-test';

interface IProps {
  form: TDraftTestForm<IPostmanValues>;
}

const UploadCollectionForm = ({form}: IProps) => {
  return (
    <div style={{display: 'grid'}}>
      <Row gutter={12}>
        <Col span={12}>
          <Form.Item name="requests" hidden>
            <Input type="hidden" />
          </Form.Item>
          <Form.Item name="variables" hidden>
            <Input type="hidden" />
          </Form.Item>
          <CollectionFileField form={form} />
          <EnvFileField form={form} />
          <SelectTestFromCollection form={form} />
        </Col>
      </Row>
      <Row gutter={12} style={{marginTop: 16}}>
        <Col span={12}>
          <RequestDetailsUrlInput />
        </Col>
        <Col span={12}>
          <Form.Item className="input-body" data-cy="body" label="Request body" name="body" style={{marginBottom: 0}}>
            <Input.TextArea placeholder="Enter request body text" />
          </Form.Item>
        </Col>
      </Row>
      <Row gutter={12}>
        <Col span={12}>
          <RequestDetailsHeadersInput />
        </Col>
        <Col span={12}>
          <RequestDetailsAuthInput form={form} />
        </Col>
      </Row>
    </div>
  );
};

export default UploadCollectionForm;
