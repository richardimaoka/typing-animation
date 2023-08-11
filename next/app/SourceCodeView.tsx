"use client";
import { useEffect, useState } from "react";

export interface Chunk {
  Content: string;
  Type: "ADD" | "DELETE" | "EQUAL";
}

interface InProgress {
  kind: "InProgress";
  currentChunk: number;
  inChunkPos: number;
  overallPos: number;
}

interface Done {
  kind: "Done";
}

type State = InProgress | Done;

interface SourceCodeViewProps {
  chunks: Chunk[];
}

const insertChar = (src: string, pos: number, char: string): string => {
  return src.slice(0, pos) + char + src.slice(pos);
};

const removeChar = (src: string, pos: number): string => {
  return src.slice(0, pos) + src.slice(pos + 1);
};

export const SourceCodeView = ({ chunks }: SourceCodeViewProps) => {
  const [sourceCode, setSourceCode] = useState(`1111
2222
3333
4444
5555`);

  const [state, setState] = useState<State>({
    kind: "InProgress",
    currentChunk: 0,
    inChunkPos: 0,
    overallPos: 0,
  });

  useEffect(() => {
    if (state.kind === "Done") {
      return;
    }

    const chunk = chunks[state.currentChunk];

    console.log(chunk.Type.padStart(5, " "), JSON.stringify(state));

    switch (chunk.Type) {
      case "EQUAL":
        const nextChunk = state.currentChunk + 1;
        if (nextChunk > chunks.length - 1) {
          setState({ kind: "Done" });
        } else {
          setState({
            kind: "InProgress",
            currentChunk: nextChunk,
            inChunkPos: 0,
            overallPos: state.overallPos + chunk.Content.length,
          });
        }
        break;
      case "ADD":
        if (state.inChunkPos === chunk.Content.length) {
          // this chunk is finished
          const nextChunk = state.currentChunk + 1;
          if (nextChunk > chunks.length - 1) {
            setState({ kind: "Done" });
          } else {
            setState({
              kind: "InProgress",
              currentChunk: nextChunk,
              inChunkPos: 0,
              overallPos: state.overallPos,
            });
          }
        } else {
          // keep processing this chunk
          setSourceCode(
            insertChar(
              sourceCode,
              state.overallPos,
              chunk.Content[state.inChunkPos]
            )
          );

          const nextChunkPos = state.inChunkPos + 1;
          const nextOverallPos = state.overallPos + 1;
          setState({
            kind: "InProgress",
            currentChunk: state.currentChunk,
            inChunkPos: nextChunkPos,
            overallPos: nextOverallPos,
          });
        }
        break;
      case "DELETE":
        if (state.inChunkPos === chunk.Content.length) {
          // this chunk is finished
          const nextChunk = state.currentChunk + 1;
          if (nextChunk > chunks.length - 1) {
            setState({ kind: "Done" });
          } else {
            setState({
              kind: "InProgress",
              currentChunk: nextChunk,
              inChunkPos: 0,
              overallPos: state.overallPos,
            });
          }
        } else {
          // keep processing this chunk
          setSourceCode(removeChar(sourceCode, state.overallPos));

          const nextChunkPos = state.inChunkPos + 1;
          const nextOverallPos = state.overallPos;
          setState({
            kind: "InProgress",
            currentChunk: state.currentChunk,
            inChunkPos: nextChunkPos,
            overallPos: nextOverallPos,
          });
        }
        break;
    }
  }, [chunks, sourceCode, state]);

  return (
    <div>
      <pre>
        <code>{sourceCode} </code>
      </pre>
      {state.kind === "InProgress" && <div>{`${JSON.stringify(state)}`}</div>}
    </div>
  );
};
