import fs from "node:fs";
import { EditOperation, SourceCodeEditor } from "./monaco/SourceCodeEditor";

interface Props {}

export async function Right(props: Props) {
  const editorText = fs.readFileSync(
    process.cwd() + "/app/components/data/command.9334a46.go.txt",
    "utf8"
  );

  const editsText = fs.readFileSync(
    process.cwd() + "/app/components/data/9334a46-51f06c7.edits.json",
    "utf8"
  );
  const edits = JSON.parse(editsText) as EditOperation[];
  console.log("edits", edits);

  return (
    <SourceCodeEditor
      editorText={editorText}
      language="go"
      typingAnimation={{
        editSequence: {
          id: "aaa",
          edits: edits,
        },
        newEditorText: "",
      }}
    />
  );
}
