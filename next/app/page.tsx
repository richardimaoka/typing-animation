import {
  getBranches,
  getCommits,
  getEdits,
  getFiles,
  getRepo,
} from "@/api/api";
import { Left } from "./components/Left";
import styles from "./page.module.css";
import { promises as fs } from "fs";
import { Right } from "./components/Right";

interface Props {
  params: { slug: string }; // for dynamic routes only
  searchParams: {
    orgname?: string | string[];
    reponame?: string | string[];
    branch?: string | string[];
    filepath?: string | string[];
  };
}

function strOrUndef(param: string | string[] | undefined): string | undefined {
  return typeof param === "string" && param !== "" ? param : undefined;
}

export default async function Page(props: Props) {
  const orgname = strOrUndef(props.searchParams.orgname);
  const reponame = strOrUndef(props.searchParams.reponame);
  const branch = strOrUndef(props.searchParams.branch);
  const filepath = strOrUndef(props.searchParams.filepath);

  const repo = await getRepo(orgname, reponame);
  const branches = await getBranches(orgname, reponame);
  const files = await getFiles(orgname, reponame, branch);
  const commits = await getCommits(orgname, reponame, branch, filepath);
  const edits = await getEdits(orgname, reponame, filepath, undefined);

  return (
    <div className={styles.component}>
      <Left
        orgname={orgname}
        reponame={reponame}
        branch={branch}
        branchSelection={branches}
        files={files}
        filepath={filepath}
        commits={commits}
        repoReady={repo.status === "ready"}
      />
      <Right />
    </div>
  );
}
