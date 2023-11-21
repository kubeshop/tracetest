import {ArrowLeftOutlined} from '@ant-design/icons';
import {Button, Divider} from 'antd';
import TestSpecActions from 'components/TestSpec/Actions';
import * as STestSpec from 'components/TestSpec/TestSpec.styled';
import {singularOrPlural} from 'utils/Common';
import {SupportedEditors} from 'constants/Editor.constants';
import {Editor} from 'components/Inputs';
import * as S from './TestSpecDetail.styled';

interface IProps {
  affectedSpans: number;
  assertionsFailed: number;
  assertionsPassed: number;
  isDeleted: boolean;
  isDraft: boolean;
  onClose(): void;
  onDelete(): void;
  onEdit(): void;
  onRevert(): void;
  selector: string;
  title: string;
}

const Header = ({
  affectedSpans,
  assertionsFailed,
  assertionsPassed,
  isDeleted,
  isDraft,
  onClose,
  onDelete,
  onEdit,
  onRevert,
  selector,
  title,
}: IProps) => {
  return (
    <>
      <S.HeaderContainer>
        <S.Row $hasGap>
          <Button icon={<ArrowLeftOutlined />} onClick={onClose} type="link" />
          <S.HeaderTitle level={2}>Test Spec Detail</S.HeaderTitle>
          <div>
            {Boolean(assertionsPassed) && (
              <STestSpec.HeaderDetail>
                <STestSpec.HeaderDot $passed />
                {assertionsPassed}
              </STestSpec.HeaderDetail>
            )}
            {Boolean(assertionsFailed) && (
              <STestSpec.HeaderDetail>
                <STestSpec.HeaderDot $passed={false} />
                {assertionsFailed}
              </STestSpec.HeaderDetail>
            )}
            <STestSpec.HeaderDetail>
              <STestSpec.HeaderSpansIcon />
              {`${affectedSpans} ${singularOrPlural('span', affectedSpans)}`}
            </STestSpec.HeaderDetail>
          </div>
        </S.Row>

        <TestSpecActions
          isDeleted={isDeleted}
          isDraft={isDraft}
          onDelete={onDelete}
          onEdit={onEdit}
          onRevert={onRevert}
        />
      </S.HeaderContainer>

      <Divider />

      {title && <S.HeaderTitle level={3}>{title}</S.HeaderTitle>}

      {selector && (
        <Editor
          type={SupportedEditors.Selector}
          editable={false}
          basicSetup={{lineNumbers: false}}
          placeholder="Selecting All Spans"
          value={selector}
        />
      )}

      <Divider />
    </>
  );
};

export default Header;
