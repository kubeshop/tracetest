import {AimOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import {useCallback} from 'react';

import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {TEST_SPEC_SNIPPETS, TSnippet} from 'constants/TestSpecs.constants';
import {isRunStateSucceeded} from 'models/TestRun.model';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import * as S from './TestSpecs.styled';

const Empty = () => {
  const {
    run: {state},
  } = useTestRun();
  const {open} = useTestSpecForm();

  const onSnippetClick = useCallback(
    (snippet: TSnippet) => {
      open({
        isEditing: false,
        selector: snippet.selector,
        defaultValues: snippet,
      });
    },
    [open]
  );

  return (
    <S.EmptyContainer data-cy="empty-test-specs">
      <S.EmptyIcon />
      <S.EmptyTitle>There are no test specs for this test</S.EmptyTitle>
      <S.EmptyText>Add a test spec, or choose from our predefined test specs:</S.EmptyText>
      <S.SnippetsContainer>
        {TEST_SPEC_SNIPPETS.map((snippet, index) => (
          // eslint-disable-next-line react/no-array-index-key
          <div key={`${snippet.selector}-${index}`}>
            <Button
              disabled={!isRunStateSucceeded(state)}
              icon={<AimOutlined />}
              onClick={() => onSnippetClick(snippet)}
              type="link"
            >
              {snippet.name}
            </Button>
          </div>
        ))}
      </S.SnippetsContainer>
    </S.EmptyContainer>
  );
};

export default Empty;
