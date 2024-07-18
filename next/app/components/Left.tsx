import styles from "./Left.module.css";

interface Props {}

export function Left(props: Props) {
  return (
    <div className={styles.component}>
      <label htmlFor="url" className={styles.label}>
        GitHub repository URL
      </label>
      <input
        id="url"
        className={styles.input}
        placeholder="https://github.com/{orgname}/{reponame}"
      />

      <label htmlFor="orgname" className={styles.label}>
        <div>GitHub organization</div>
        <div>(user)</div>
      </label>
      <input id="orgname" className={styles.input} placeholder="{orgname}" />

      <label htmlFor="branch" className={styles.label}>
        branch
      </label>
      <input id="branch" className={styles.input} placeholder="{branch}" />

      <label htmlFor="filepath" className={styles.label}>
        file path
      </label>
      <input id="filepath" className={styles.input} placeholder="{filePath}" />

      <label htmlFor="commits" className={styles.label + " " + styles.top}>
        commits
      </label>
      <div id="commits" className={styles.commits} />
    </div>
  );
}
