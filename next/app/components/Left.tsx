"use client";

import styles from "./Left.module.css";
import { useRouter, useSearchParams } from "next/navigation";

interface Props {}

export function Left(props: Props) {
  const router = useRouter();

  const searchParams = useSearchParams();
  const orgnameParam = searchParams.get("orgname");
  const reponameParam = searchParams.get("reponame");
  const branchParam = searchParams.get("branch");
  const filepathParam = searchParams.get("filepath");

  function newPath(
    newOrg: string | null,
    newRepo: string | null,
    newBranch: string | null,
    newFilepath: string | null
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
      params.push("filepath=" + newFilepath);
    }

    if (params.length === 0) {
      return "/";
    } else {
      return "/?" + params.join("&");
    }
  }

  function onReponameChange(newRepo: string) {
    const href = newPath(orgnameParam, newRepo, branchParam, filepathParam);
    router.push(href);
  }

  function onOrgnameChange(newOrg: string) {
    const href = newPath(newOrg, reponameParam, branchParam, filepathParam);
    router.push(href);
  }

  function onFilePathChange(newFilePath: string) {
    const href = newPath(orgnameParam, reponameParam, branchParam, newFilePath);
    router.push(href);
  }

  function onBranchChange(newBranch: string) {
    const href = newPath(orgnameParam, reponameParam, newBranch, filepathParam);
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

      <label className={styles.label + " " + styles.grey}>GitHub URL</label>
      <div className={styles.grey}>
        https://github.com/{orgnameParam || "{orgname}"}/
        {reponameParam || "{reponame}"}
      </div>

      <label htmlFor="branch" className={styles.label}>
        branch
      </label>
      <input
        id="branch"
        className={styles.input}
        placeholder="main"
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

      <label htmlFor="commits" className={styles.label + " " + styles.top}>
        commits
      </label>
      <div id="commits" className={styles.commits} />
    </div>
  );
}
