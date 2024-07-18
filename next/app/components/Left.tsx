"use client";

import { useState } from "react";
import styles from "./Left.module.css";
import { useRouter, useSearchParams } from "next/navigation";

interface Props {}

export function Left(props: Props) {
  const router = useRouter();

  const searchParams = useSearchParams();
  const orgname = searchParams.get("orgname");
  const reponame = searchParams.get("reponame");

  function newURL(newOrg: string | null, newRepo: string | null): string {
    let params = [];

    if (newOrg && newOrg !== "") {
      params.push("orgname=" + newOrg);
    }
    if (newRepo && newRepo !== "") {
      params.push("reponame=" + newRepo);
    }

    if (params.length === 0) {
      return "/";
    } else {
      return "/?" + params.join("&");
    }
  }

  function onReponameChange(newRepo: string) {}

  function onOrgnameChange(newOrg: string) {
    const url = newURL(newOrg, reponame);
    router.push(url);
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
      <input id="reponame" className={styles.input} placeholder="{reponame}" />

      <label className={styles.label + " " + styles.grey}>GitHub URL</label>
      <div className={styles.grey}>
        https://github.com/{orgname || "{orgname}"}/{reponame || "{reponame}"}
      </div>

      <label htmlFor="branch" className={styles.label}>
        branch
      </label>
      <input
        id="branch"
        className={styles.input}
        defaultValue="main"
        placeholder="{branch}"
      />

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