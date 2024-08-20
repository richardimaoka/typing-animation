import fs from "node:fs";
import { EditOperation, SourceCodeEditor } from "./monaco/SourceCodeEditor";

interface Props {
  editorText?: string;
}

export async function Right(props: Props) {
  if (!props.editorText) {
    return <></>;
  }

  return (
    <SourceCodeEditor
      editorText={props.editorText}
      language="go"
      // typingAnimation={{
      //   editSequence: {
      //     id: "aaa",
      //     edits: edits,
      //   },
      //   newEditorText: "",
      // }}
    />
  );
}
