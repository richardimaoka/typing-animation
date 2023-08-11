import { Chunk, SourceCodeView } from "./SourceCodeView";

export default function Home() {
  const chunks: Chunk[] = [
    {
      Content: "1111\n",
      Type: "EQUAL",
    },
    {
      Content: "2222\n",
      Type: "DELETE",
    },
    {
      Content: "222222\n",
      Type: "ADD",
    },
    {
      Content: "",
      Type: "EQUAL",
    },
    {
      Content: "3333\n",
      Type: "DELETE",
    },
    {
      Content: "333333\n",
      Type: "ADD",
    },
    {
      Content: "4444\n",
      Type: "EQUAL",
    },
    {
      Content: "5555\n",
      Type: "DELETE",
    },
    {
      Content: "5555555\n",
      Type: "ADD",
    },
  ];

  return (
    <main>
      <SourceCodeView chunks={chunks} />
    </main>
  );
}
