import { Chunk, SourceCodeView } from "./SourceCodeView";

export default function Home() {
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
    <main>
      <SourceCodeView chunks={chunks} />
    </main>
  );
}
