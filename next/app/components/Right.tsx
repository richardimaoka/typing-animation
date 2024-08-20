import fs from "node:fs";
import { EditOperation, SourceCodeEditor } from "./monaco/SourceCodeEditor";

interface Props {
  editorText?: string;
  edits?: EditOperation[];
}

export async function Right(props: Props) {
  if (!props.editorText) {
    return <></>;
  }

  return (
    <SourceCodeEditor
      editorText={props.editorText}
      language="go"
      typingAnimation={
        props.edits
          ? {
              editSequence: {
                id: "aaa",
                edits: props.edits,
              },
              newEditorText: "",
            }
          : undefined
      }
    />
  );
}
