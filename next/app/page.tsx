import { getBranches, getCommits, getFiles, getRepo } from "@/api/api";
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

function toParamString(
  param: string | string[] | undefined
): string | undefined {
  return typeof param === "string" && param !== "" ? param : undefined;
}

export default async function Page(props: Props) {
  const orgname = toParamString(props.searchParams.orgname);
  const reponame = toParamString(props.searchParams.reponame);
  const branch = toParamString(props.searchParams.branch) || "main";
  const filepath = toParamString(props.searchParams.filepath) || "main";

  const repo = await getRepo(orgname, reponame);
  if (!repo) {
    return <></>;
  }

  const branches = await getBranches(orgname, reponame);
  const files = await getFiles(orgname, reponame, branch);
  const commits = await getCommits(orgname, reponame, branch, filepath);

  return (
    <div className={styles.component}>
      <Left
        orgname={orgname}
        reponame={reponame}
        branch={branch}
        files={files}
        filepath={filepath}
        commits={commits}
      />
      <Right />
    </div>
  );
}
