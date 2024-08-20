import fs from "node:fs";
import { EditOperation, SourceCodeEditor } from "./monaco/SourceCodeEditor";

interface Props {
  editorText?: string;
  edits?: EditOperation[];
  commit?: string;
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
                id: props.commit || "aaa",
                edits: props.edits,
              },
              newEditorText: "",
            }
          : undefined
      }
    />
  );
}
