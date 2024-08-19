import { CommitData } from "@/api/types";
import styles from "./Left.module.css";
import { OrgNameField } from "./fields/OrgNameField";
import { RepositoryNameField } from "./fields/RepositoryNameField";
import { GitHubURLDisplay } from "./fields/GitHubURLDisplay";
import { BranchSelection } from "./fields/BranchSelection";
import { FilePathSelection } from "./fields/FilePathSelection";

interface Props {
  commits?: CommitData[];
  files?: string[];

  orgname?: string;
  reponame?: string;
  branch?: string;
  branchSelection?: string[];
  filepath?: string;
  repoReady?: boolean;
}

export function Left(props: Props) {
  return (
    <div className={styles.component}>
      <form className={styles.repo}>
        <OrgNameField orgname={props.orgname} />
        <RepositoryNameField reponame={props.reponame} />
        <GitHubURLDisplay orgname={props.orgname} reponame={props.reponame} />
      </form>

      {props.repoReady && (
        <div className={styles.lower}>
          <BranchSelection
            branch={props.branch}
            brancheSelection={props.branchSelection}
          />
          <FilePathSelection
            filepath={props.filepath}
            fileSelection={props.files}
          />
          <label className={styles.label + " " + styles.top}>commits</label>
          <fieldset className={styles.commits}>
            {props.commits &&
              props.commits.map((c) => (
                <div key={c.hash}>
                  <input type="radio" id={c.hash} name="commit" />
                  <label htmlFor={c.hash}>
                    <span className={styles.hash}>{c.shortHash}</span>
                    <span className={styles.message}>{c.shortMessage}</span>
                  </label>
                </div>
              ))}
          </fieldset>
        </div>
      )}
    </div>
  );
}
