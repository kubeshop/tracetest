import {EditorState} from '@codemirror/state';
import {Editor} from 'components/Inputs';
import {SupportedEditors} from 'constants/Editor.constants';
import {IEditorProps} from '../Editor/Editor';

const extensions = [EditorState.transactionFilter.of(tr => (tr.newDoc.lines > 1 ? [] : tr))];

const SingleLine = (props: IEditorProps) => (
  <Editor extensions={extensions} type={SupportedEditors.Interpolation} {...props} />
);

export default SingleLine;
