import { CommitData } from "./types";
import fs from "node:fs";

type RepoStatus = {
  status: "ready" | "in progress" | "wrong name" | "error";
};

export async function getRepo(
  orgname: string | undefined,
  reponame: string | undefined
): Promise<RepoStatus> {
  if (!orgname) {
    return { status: "wrong name" };
  }

  if (!reponame) {
    return { status: "wrong name" };
  }

  // For fetch API, try-catch is easier than await-catch,
  // because await-catch requires you to craft a `Response` variable
  let response: Response; /* let, instead of const, is easier for try-catch, but make sure you NEVER re-assign to the let variable */
  try {
    response = await fetch(`http://localhost:8080/${orgname}/${reponame}`, {
      cache: "no-store",
    });
  } catch (error) {
    // Only network errors reach here.
    // (i.e.) HTTP response with an error status (e.g. 500) does NOT come into this catch block
    console.error("getRepo failed (network error)", error);
    // by using try-catch, we are able to directly return from function
    return { status: "error" };
  }

  if (!response.ok) {
    return { status: "error" };
  }

  const jsonData = (await response.json().catch(function (error) {
    if (error instanceof SyntaxError) {
      console.log("There was a SyntaxError in returned JSON", error);
    } else {
      console.log("There was an upon parsing JSON response", error);
    }
    return { status: "error" };
  })) as RepoStatus;

  console.log("getRepo successful", jsonData);
  return jsonData;
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
