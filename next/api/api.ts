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

type FilesData = {
  orgname: string;
  repo: string;
  files: string[];
};

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

  // For fetch API, try-catch is easier than await-catch,
  // because await-catch requires you to craft a `Response` variable
  let response: Response; /* let, instead of const, is easier for try-catch, but make sure you NEVER re-assign to the let variable */
  try {
    response = await fetch(
      `http://localhost:8080/${orgname}/${reponame}/files`,
      {
        cache: "no-store",
      }
    );
  } catch (error) {
    // Only network errors reach here.
    // (i.e.) HTTP response with an error status (e.g. 500) does NOT come into this catch block
    console.error("getRepo failed (network error)", error);

    // by using try-catch, we are able to directly return from function
    return undefined;
  }

  if (!response.ok) {
    return undefined;
  }

  const jsonData = (await response.json().catch(function (error) {
    if (error instanceof SyntaxError) {
      console.log("There was a SyntaxError in returned JSON", error);
    } else {
      console.log("There was an upon parsing JSON response", error);
    }
    return { status: "error" };
  })) as FilesData;

  console.log("getFiles successful", jsonData);

  return jsonData.files;
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

  // For fetch API, try-catch is easier than await-catch,
  // because await-catch requires you to craft a `Response` variable
  let response: Response; /* let, instead of const, is easier for try-catch, but make sure you NEVER re-assign to the let variable */
  try {
    const path = decodeURIComponent(filepath);
    response = await fetch(
      `http://localhost:8080/${orgname}/${reponame}/files/${path}`,
      {
        cache: "no-store",
      }
    );
  } catch (error) {
    // Only network errors reach here.
    // (i.e.) HTTP response with an error status (e.g. 500) does NOT come into this catch block
    console.error("getCommits failed (network error)", error);

    // by using try-catch, we are able to directly return from function
    return undefined;
  }

  if (!response.ok) {
    return undefined;
  }

  const jsonData = (await response.json().catch(function (error) {
    if (error instanceof SyntaxError) {
      console.log("There was a SyntaxError in returned JSON", error);
    } else {
      console.log("There was an upon parsing JSON response", error);
    }
    return { status: "error" };
  })) as FilesData;

  console.log("getFiles successful", jsonData);

  return jsonData.files;
}
