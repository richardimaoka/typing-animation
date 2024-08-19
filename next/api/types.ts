import { editor } from "monaco-editor";

export type CommitData = {
  hash: string;
  shortHash: string;
  shortMessage: string;
};

export type FileData = {
  commits: CommitData[];
  contents: string;
  edits: editor.ISingleEditOperation[];
};
