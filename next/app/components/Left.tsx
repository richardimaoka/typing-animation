"use client";

import styles from "./Left.module.css";
import { useRouter, useSearchParams } from "next/navigation";

type CommitData = {
  hash: string;
  message: string;
};

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

  function onReponameChange(newRepo: string) {
    const href = newPath(props.orgname, newRepo, props.branch, props.filepath);
    router.push(href);
  }

  function onOrgnameChange(newOrg: string) {
    const href = newPath(newOrg, props.reponame, props.branch, props.filepath);
    router.push(href);
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
      {/* <label htmlFor="url" className={styles.label}>
        GitHub repository URL
      </label>
      <input
        id="url"
        className={styles.input}
        placeholder="https://github.com/{orgname}/{reponame}"
        onBlur={(e) => {
          console.log("blur", e.target.value, e);
        }}
        onKeyDown={(e) => {
          console.log("keydown", e);
        }}
      /> */}

      <label htmlFor="orgname" className={styles.label}>
        <div>GitHub organization</div>
        <div>(user)</div>
      </label>
      <input
        id="orgname"
        className={styles.input}
        defaultValue={props.orgname}
        placeholder="{orgname}"
        onBlur={(e) => {
          onOrgnameChange(e.target.value);
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            onOrgnameChange(e.currentTarget.value);
          }
        }}
      />

      <label htmlFor="reponame" className={styles.label}>
        GitHub repository
      </label>
      <input
        id="reponame"
        className={styles.input}
        defaultValue={props.reponame}
        placeholder="{reponame}"
        onBlur={(e) => {
          onReponameChange(e.target.value);
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            onReponameChange(e.currentTarget.value);
          }
        }}
      />

      <label className={styles.label}>GitHub URL</label>
      <div className={styles.grey}>
        https://github.com/{props.orgname || "{orgname}"}/
        {props.reponame || "{reponame}"}
      </div>

      {/* Empty space in CSS grid */}
      <label />
      <div></div>

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
      <input
        id="filepath"
        className={styles.input}
        placeholder="{filepath}"
        onBlur={(e) => {
          onFilePathChange(e.target.value);
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter") {
            onFilePathChange(e.currentTarget.value);
          }
        }}
      />

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
  );
}
