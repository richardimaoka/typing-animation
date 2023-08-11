import { useEffect, useState } from "react";

export const SourceCodeView = () => {
  const [sourceCode, setSourceCode] = useState(`1111
2222
3333
4444
5555`);

  interface Chunk {
    Content: string;
    Type: "ADD" | "DELETE" | "EQUAL";
  }
  const chunks: Chunk[] = [
    {
      Content: "1111\n2222",
      Type: "EQUAL",
    },
    {
      Content: "22",
      Type: "ADD",
    },
    {
      Content: "\n",
      Type: "EQUAL",
    },
    {
      Content: "33",
      Type: "ADD",
    },
    {
      Content: "3333\n4444\n",
      Type: "EQUAL",
    },
    {
      Content: "555",
      Type: "ADD",
    },
    {
      Content: "5555\n",
      Type: "EQUAL",
    },
  ];

  return (
    <pre>
      <code>{sourceCode} </code>
    </pre>
  );
};
