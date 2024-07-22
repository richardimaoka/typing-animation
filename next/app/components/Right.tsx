import { SourceCodeEditor } from "./monaco/SourceCodeEditor";
import styles from "./Right.module.css";
import fs from "node:fs";

interface Props {}

export async function Right(props: Props) {
  const editorText = fs.readFileSync(
    process.cwd() + "/app/components/data/command.9334a46.go.txt",
    "utf8"
  );

  const editsText = fs.readFileSync(
    process.cwd() + "/app/components/data/51f06c7-9334a46.edits.json",
    "utf8"
  );

  const edits = JSON.parse(editsText);

  return (
    <SourceCodeEditor
      editorText={editorText}
      language="go"
      typingAnimation={{
        editSequence: edits,
        newEditorText: "",
      }}
    />
  );
}
