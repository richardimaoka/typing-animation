import styles from "./GitHubURLDisplay.module.css";

interface Props {
  orgname?: string | undefined;
  reponame?: string | undefined;
  repoReady?: boolean;
}

export function GitHubURLDisplay(props: Props) {
  return (
    <>
      <label className={styles.label}>GitHub URL</label>
      <div className={styles.grey}>
        https://github.com/{props.orgname || "{orgname}"}/
        {props.reponame || "{reponame}"}
        {props.repoReady ? (
          <span>ready</span>
        ) : (
          <button className={styles.button}>load</button>
        )}
      </div>
    </>
  );
}
