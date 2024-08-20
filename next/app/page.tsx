import {
  getBranches,
  getCommits,
  getEdits,
  getFileContents,
  getFiles,
  getRepo,
} from "@/api/api";
import { Left } from "./components/Left";
import styles from "./page.module.css";
import { Right } from "./components/Right";

interface Props {
  params: { slug: string }; // for dynamic routes only
  searchParams: {
    orgname?: string | string[];
    reponame?: string | string[];
    branch?: string | string[];
    filepath?: string | string[];
    commit?: string | string[];
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
  const commit = strOrUndef(props.searchParams.commit);

  const repo = await getRepo(orgname, reponame);
  const branches = await getBranches(orgname, reponame);
  const files = await getFiles(orgname, reponame, branch);
  const commits = await getCommits(orgname, reponame, branch, filepath);

  const contents = await getFileContents(orgname, reponame, filepath, commit);
  const edits = await getEdits(orgname, reponame, filepath, commit);

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
        commit={commit}
      />
      <Right editorText={contents} edits={edits} commit={commit} />
    </div>
  );
}
