import { CommitData } from "./types";

export async function getFiles(
  orgname: string | undefined,
  reponame: string | undefined,
  branch: string | undefined
): Promise<string[] | undefined> {
  if (!orgname) {
    return undefined;
  }
  if (!reponame) {
    return undefined;
  }
  if (!branch) {
    return undefined;
  }

  return ["a.go", "b.go", "c.go"];
}

export async function getCommits(
  orgname: string | undefined,
  reponame: string | undefined,
  filepath: string | undefined
): Promise<CommitData[] | undefined> {
  if (!orgname) {
    return undefined;
  }
  if (!reponame) {
    return undefined;
  }
  if (!filepath) {
    return undefined;
  }

  return [
    { hash: "sad65sd234f", message: "this is the first commit" },
    { hash: "asd897sdf87", message: "this is the second commit" },
    { hash: "890908sdr49", message: "this is the third commit" },
    { hash: "78searoa36a", message: "this is the fourth commit" },
    { hash: "905sfdjo439", message: "this is the fifth commit" },
  ];
}
