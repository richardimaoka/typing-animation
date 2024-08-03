"use client";

import { CommitData } from "@/api/types";
import { useRouter, useSearchParams } from "next/navigation";
import styles from "./Left.module.css";
import { OrgNameField } from "./fields/OrgNameField";
import { RepositoryNameField } from "./fields/RepositoryNameField";
import { GitHubURLDisplay } from "./fields/GitHubURLDisplay";

interface Props {
  commits?: CommitData[];
  files?: string[];

  orgname?: string;
  reponame?: string;
  branch: string;
  filepath?: string;
}

export function Left(props: Props) {
  const router = useRouter();

  function newPath(
    newOrg: string | undefined,
    newRepo: string | undefined,
    newBranch: string | undefined,
    newFilepath: string | undefined
  ): string {
    let params = [];

    if (newOrg && newOrg !== "") {
      params.push("orgname=" + newOrg);
    }

    if (newRepo && newRepo !== "") {
      params.push("reponame=" + newRepo);
    }

    if (newBranch && newBranch !== "") {
      params.push("branch=" + newBranch);
    }

    if (newFilepath && newFilepath !== "") {
      params.push("filepath=" + encodeURIComponent(newFilepath));
    }

    if (params.length === 0) {
      return "/";
    } else {
      return "/?" + params.join("&");
    }
  }

  function onFilePathChange(newFilePath: string) {
    const href = newPath(
      props.orgname,
      props.reponame,
      props.branch,
      newFilePath
    );
    router.push(href);
  }

  function onBranchChange(newBranch: string) {
    const href = newPath(
      props.orgname,
      props.reponame,
      newBranch,
      props.filepath
    );
    router.push(href);
  }

  return (
    <div className={styles.component}>
      <form className={styles.repo}>
        <OrgNameField orgname={props.orgname} />
        <RepositoryNameField reponame={props.reponame} />
        <GitHubURLDisplay orgname={props.orgname} reponame={props.reponame} />
      </form>

      <div className={styles.lower}>
        <label htmlFor="branch" className={styles.label}>
          branch
        </label>
        <input
          id="branch"
          className={styles.input}
          placeholder="main"
          defaultValue={props.branch}
          onBlur={(e) => {
            onBranchChange(e.target.value);
          }}
          onKeyDown={(e) => {
            if (e.key === "Enter") {
              onBranchChange(e.currentTarget.value);
            }
          }}
        />
        <label htmlFor="filepath" className={styles.label}>
          file path
        </label>
        <select
          id="filepath"
          onChange={(e) => {
            onFilePathChange(e.target.value);
          }}
        >
          {props.files ? (
            props.files.map((fpath) => (
              <option key={fpath} value={fpath}>
                {fpath}
              </option>
            ))
          ) : (
            <option disabled>no option is available</option>
          )}
        </select>
        <label className={styles.label + " " + styles.top}>commits</label>
        <fieldset className={styles.commits}>
          {props.commits &&
            props.commits.map((c) => (
              <div key={c.hash}>
                <input type="radio" id={c.hash} name="commit" />
                <label htmlFor={c.hash}>
                  <span className={styles.hash}>{c.hash}</span>
                  <span className={styles.message}>{c.message}</span>
                </label>
              </div>
            ))}
        </fieldset>
      </div>
    </div>
  );
}
