import { SourceCodeEditor } from "./monaco/SourceCodeEditor";
import styles from "./Right.module.css";

interface Props {}

export function Right(props: Props) {
  return <SourceCodeEditor editorText="aa" language="go" />;
}
