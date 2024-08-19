import { CommitData } from "@/api/types";
import styles from "./Left.module.css";
import { OrgNameField } from "./fields/OrgNameField";
import { RepositoryNameField } from "./fields/RepositoryNameField";
import { GitHubURLDisplay } from "./fields/GitHubURLDisplay";
import { BranchSelection } from "./fields/BranchSelection";
import { FilePathSelection } from "./fields/FilePathSelection";
import { CommitSelection } from "./fields/CommitSelection";

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
          <CommitSelection commit={""} commitSelection={props.commits} />
        </div>
      )}
    </div>
  );
}
