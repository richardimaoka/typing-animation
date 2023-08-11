"use client";
import { useEffect, useState } from "react";

export interface Chunk {
  Content: string;
  Type: "ADD" | "DELETE" | "EQUAL";
}

interface Initializing {
  kind: "Initializing";
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

type State = Initializing | InProgress | Done;

const transition = (
  chunks: Chunk[],
  state: InProgress,
  sourceCode: string
): [State, string] => {
  const chunk = chunks[state.currentChunk];

  switch (chunk.Type) {
    case "EQUAL":
      const nextChunk = state.currentChunk + 1;
      if (nextChunk > chunks.length - 1) {
        return [{ kind: "Done" }, sourceCode];
      } else {
        return [
          {
            kind: "InProgress",
            currentChunk: nextChunk,
            inChunkPos: 0,
            overallPos: state.overallPos + chunk.Content.length,
          },
          sourceCode,
        ];
      }
    case "ADD":
      if (state.inChunkPos === chunk.Content.length) {
        // this chunk is finished
        const nextChunk = state.currentChunk + 1;
        if (nextChunk > chunks.length - 1) {
          return [{ kind: "Done" }, sourceCode];
        } else {
          return [
            {
              kind: "InProgress",
              currentChunk: nextChunk,
              inChunkPos: 0,
              overallPos: state.overallPos,
            },
            sourceCode,
          ];
        }
      } else {
        // keep processing this chunk
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
        ];
      }
    case "DELETE":
      if (state.inChunkPos === chunk.Content.length) {
        // this chunk is finished
        const nextChunk = state.currentChunk + 1;
        if (nextChunk > chunks.length - 1) {
          return [{ kind: "Done" }, sourceCode];
        } else {
          return [
            {
              kind: "InProgress",
              currentChunk: nextChunk,
              inChunkPos: 0,
              overallPos: state.overallPos,
            },
            sourceCode,
          ];
        }
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
  const [state, setState] = useState<State>({
    kind: "Initializing",
  });

  useEffect(() => {
    if (state.kind === "Initializing") {
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
      const [newState, newSourceCode] = transition(chunks, state, src);
      setTimeout(() => {
        setState(newState);
        setSourceCode(newSourceCode);
      }, 100);
    }
  }, [chunks, sourceCode, src, state]);

  return (
    <div>
      <pre>
        <code>{src}</code>
      </pre>
      <div>{`${JSON.stringify(state)}`}</div>
    </div>
  );
};
