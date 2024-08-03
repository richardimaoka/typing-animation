import { CommitData } from "./types";
import fs from "node:fs";

type RepoStatus = {};

export async function getRepo(
  orgname: string | undefined,
  reponame: string | undefined
): Promise<string[] | undefined> {
  if (!orgname) {
    return undefined;
  }

  if (!reponame) {
    return undefined;
  }

  if (orgname === "spf13" && reponame === "cobra") {
    return [""];
  }

  return undefined; //[]; //"a.go", "b.go", "c.go"];
}

export async function getBranches(
  orgname: string | undefined,
  reponame: string | undefined
): Promise<string[] | undefined> {
  if (!orgname) {
    return undefined;
  }

  if (!reponame) {
    return undefined;
  }

  if (orgname === "spf13" && reponame === "cobra") {
    return ["main"];
  }

  return undefined; //[]; //"a.go", "b.go", "c.go"];
}

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

  if (orgname === "spf13" && reponame === "cobra") {
    const fileContents = fs.readFileSync(
      process.cwd() + "/api/files.json",
      "utf8"
    );
    const files = JSON.parse(fileContents) as string[];

    return files;
  }

  return undefined; //[]; //"a.go", "b.go", "c.go"];
}

export async function getCommits(
  orgname: string | undefined,
  reponame: string | undefined,
  branch: string | undefined,
  filepath: string | undefined
): Promise<CommitData[] | undefined> {
  if (!orgname) {
    return undefined;
  }

  if (!reponame) {
    return undefined;
  }

  if (!branch) {
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
