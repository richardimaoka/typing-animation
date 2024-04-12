"use client";

import { useEffect, useState } from "react";
import { Chunk, SourceCodeView } from "./SourceCodeView";

interface State {
  repository: string;
  prevCommit: string;
  currentCommit: string;
  filepath: string;
}

export default function Home() {
  const [state, setState] = useState<State>({
    repository: "typing-animation",
    prevCommit: "7d22d71cca1c9ba3c280a78d7311613e785eb5ee",
    currentCommit: "99d5c1c214319cc7b0da69e692c24f639095bc1a",
    filepath: "next/app/SourceCodeView.tsx",
  });

  const chunks: Chunk[] = [
    {
      Content:
        '\u003c!DOCTYPE html\u003e\n\u003chtml lang="en"\u003e\n  \u003chead\u003e\n    \u003cmeta charset="UTF-8" /\u003e\n    \u003cmeta name="viewport" content="width=device-width, initial-scale=1.0" /\u003e\n    \u003ctitle\u003eDocument\u003c/title\u003e\n    \u003cscript src="https://accounts.google.com/gsi/client" async defer\u003e\u003c/script\u003e\n  \u003c/head\u003e\n  \u003cbody\u003e\n    this is a blank document\n  \u003c/body\u003e\n  \u003cscript\u003e\n    window.onload = function () {\n',
      Type: "EQUAL",
    },
    {
      Content: '      console.log("loaded");\n',
      Type: "ADD",
    },
    {
      Content:
        '      google.accounts.id.initialize({\n        client_id:\n          "13173511749-e9dutacu8tmq9f8ro1bt9dh74ajqb700.apps.googleusercontent.com",\n        callback: handleCredentialResponse,\n      });\n      google.accounts.id.prompt();\n',
      Type: "EQUAL",
    },
    {
      Content:
        '    };\n\n    handleCredentialResponse = function () {\n      console.log("handleCredentialResponse called");\n',
      Type: "ADD",
    },
    {
      Content: "    };\n  \u003c/script\u003e\n\u003c/html\u003e\n",
      Type: "EQUAL",
    },
  ];

  const sourceCode = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Document</title>
    <script src="https://accounts.google.com/gsi/client" async defer></script>
  </head>
  <body>
    this is a blank document
  </body>
  <script>
    window.onload = function () {
      google.accounts.id.initialize({
        client_id:
          "13173511749-e9dutacu8tmq9f8ro1bt9dh74ajqb700.apps.googleusercontent.com",
        callback: handleCredentialResponse,
      });
      google.accounts.id.prompt();
    };
  </script>
</html>`;

  useEffect(() => {
    const a = fetch(
      `https://raw.githubusercontent.com/richardimaoka/${state.repository}/${state.prevCommit}/${state.filepath}`
    ).then((response) => response.text());
    const b = fetch(
      `https://raw.githubusercontent.com/richardimaoka/${state.repository}/${state.currentCommit}/${state.filepath}`
    ).then((response) => response.text());

    Promise.all([a, b]).then((c) => c);
  });

  return (
    <main>
      <div>
        <input
          type="text"
          placeholder="repository"
          value={state.repository}
          onChange={(e) =>
            setState({ ...state, repository: e.currentTarget.value })
          }
        />
      </div>
      <div>
        <input
          type="text"
          placeholder="prev commit"
          value={state.prevCommit}
          onChange={(e) =>
            setState({ ...state, prevCommit: e.currentTarget.value })
          }
        />
      </div>
      <div>
        <input
          type="text"
          placeholder="current commit"
          value={state.currentCommit}
          onChange={(e) =>
            setState({ ...state, currentCommit: e.currentTarget.value })
          }
        />
      </div>
      <div>
        <input
          type="text"
          placeholder="filepath"
          value={state.filepath}
          onChange={(e) =>
            setState({ ...state, filepath: e.currentTarget.value })
          }
        />
      </div>
      {<SourceCodeView sourceCode={sourceCode} chunks={chunks} />}
    </main>
  );
}
