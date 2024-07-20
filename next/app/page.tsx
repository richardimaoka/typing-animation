import { getCommits, getFiles } from "@/api/api";
import { Left } from "./components/Left";
import styles from "./page.module.css";
import { promises as fs } from "fs";

interface Props {
  params: { slug: string }; // for dynamic routes only
  searchParams: {
    orgname?: string | string[];
    reponame?: string | string[];
    branch?: string | string[];
    filepath?: string | string[];
  };
}

function retrieveParam(
  param: string | string[] | undefined
): string | undefined {
  return typeof param === "string" && param !== "" ? param : undefined;
}

export default async function Page(props: Props) {
  const orgname = retrieveParam(props.searchParams.orgname);
  const reponame = retrieveParam(props.searchParams.reponame);
  const branch = retrieveParam(props.searchParams.branch) || "main";
  const filepath = retrieveParam(props.searchParams.filepath) || "main";

  const files = await getFiles(orgname, reponame, branch);
  const commits = await getCommits(orgname, reponame, branch);

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
      <div>right</div>
    </div>
  );
}
