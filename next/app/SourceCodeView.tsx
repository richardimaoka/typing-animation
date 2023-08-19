"use client";
import { useEffect, useState } from "react";

export interface Chunk {
  Content: string;
  Type: "ADD" | "DELETE" | "EQUAL";
}

interface Init {
  kind: "Init";
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

type State = Init | InProgress | Done;

const nextChunkState = (
  overallPos: number,
  currentChunk: number,
  chunks: Chunk[]
): InProgress | Done => {
  const chunk = chunks[currentChunk];
  const nextChunk = currentChunk + 1;

  const nextOverallPos =
    chunk.Type === "EQUAL"
      ? overallPos + chunk.Content.length // if "EQAUL", chunk.Content is skipped and offset the overallPos
      : overallPos; // if "ADD" or "DELETE", state.overallPos should already be set to chunk's last position

  if (nextChunk > chunks.length - 1) {
    return { kind: "Done" };
  } else {
    return {
      kind: "InProgress", //ReadyForChunk
      currentChunk: nextChunk,
      inChunkPos: 0,
      overallPos: nextOverallPos,
    };
  }
};

const transition = (
  chunks: Chunk[],
  state: InProgress,
  sourceCode: string
): [State, string, number] => {
  const transitionMilliSeconds = 20;
  const chunk = chunks[state.currentChunk];

  switch (chunk.Type) {
    case "EQUAL":
      return [
        nextChunkState(state.overallPos, state.currentChunk, chunks), // skip to the next chunk
        sourceCode,
        transitionMilliSeconds,
      ];
    case "ADD":
      if (state.inChunkPos === chunk.Content.length) {
        // this chunk is finished
        return [
          nextChunkState(state.overallPos, state.currentChunk, chunks),
          sourceCode,
          transitionMilliSeconds,
        ];
      } else {
        // keep processing this chunk
        const nextNewLinePos = chunk.Content.indexOf("\n");
        const nextChunkPos = state.inChunkPos + 1;
        const nextOverallPos = state.overallPos + 1;
        return [
          {
            kind: "InProgress",
            currentChunk: state.currentChunk,
            inChunkPos: nextChunkPos,
            overallPos: nextOverallPos,
          },
          insertChar(
            sourceCode,
            state.overallPos,
            chunk.Content[state.inChunkPos]
          ),
          transitionMilliSeconds,
        ];
      }
    case "DELETE":
      if (state.inChunkPos === chunk.Content.length) {
        return [
          nextChunkState(state.overallPos, state.currentChunk, chunks),
          sourceCode,
          transitionMilliSeconds,
        ];
      } else {
        // keep processing this chunk
        const nextChunkPos = state.inChunkPos + 1;
        const nextOverallPos = state.overallPos;
        return [
          {
            kind: "InProgress",
            currentChunk: state.currentChunk,
            inChunkPos: nextChunkPos,
            overallPos: nextOverallPos,
          },
          removeChar(sourceCode, state.overallPos),
          transitionMilliSeconds,
        ];
      }
  }
};

const insertChar = (src: string, pos: number, char: string): string => {
  return src.slice(0, pos) + char + src.slice(pos);
};

const removeChar = (src: string, pos: number): string => {
  return src.slice(0, pos) + src.slice(pos + 1);
};
interface SourceCodeViewProps {
  sourceCode: string;
  chunks: Chunk[];
}

export const SourceCodeView = ({ sourceCode, chunks }: SourceCodeViewProps) => {
  const [src, setSourceCode] = useState("");
  const [state, setState] = useState<State>({ kind: "Init" });

  useEffect(() => {
    if (state.kind === "Init") {
      setSourceCode(sourceCode);
      setState({
        kind: "InProgress",
        currentChunk: 0,
        inChunkPos: 0,
        overallPos: 0,
      });
    } else if (state.kind === "Done") {
      return;
    } else {
      const [newState, newSourceCode, transitionMilliSeconds] = transition(
        chunks,
        state,
        src
      );
      setTimeout(() => {
        setState(newState);
        setSourceCode(newSourceCode);
      }, transitionMilliSeconds);
    }
  }, [chunks, sourceCode, src, state]);

  return (
    <div>
      <pre>
        <code>{src}</code>
      </pre>
      {/* <div>{`${JSON.stringify(state)}`}</div> */}
    </div>
  );
};
