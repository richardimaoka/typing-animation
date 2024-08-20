"use client";

import { CommitData } from "@/api/types";
import styles from "./CommitSelection.module.css";
import { useRouter, usePathname, useSearchParams } from "next/navigation";

interface Props {
  commit?: string;
  commitSelection?: CommitData[];
}

export function CommitSelection(props: Props) {
  // https://github.com/vercel/next.js/discussions/47583
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const updateSearchParams = (newValue: string) => {
    // now you got a read/write object
    const current = new URLSearchParams(Array.from(searchParams.entries())); // -> has to use this form

    // update as necessary
    const paramName = "commit";
    if (!newValue) {
      current.delete(paramName);
    } else {
      current.set(paramName, newValue);
    }

    // cast to string
    const search = current.toString();
    // or const query = `${'?'.repeat(search.length && 1)}${search}`;
    const query = search ? `?${search}` : "";

    router.push(`${pathname}${query}`);
  };

  console.log("hashes:", props.commit);

  return (
    <>
      <label className={styles.label + " " + styles.top}>commits</label>
      <fieldset className={styles.commits}>
        {props.commitSelection && props.commitSelection.length > 0 ? (
          props.commitSelection.map((c) => (
            <div key={c.hash}>
              <input
                type="radio"
                id={c.hash}
                name="commit"
                value={c.shortHash}
                onChange={(e) => {
                  updateSearchParams(e.target.value.trim());
                }}
                checked={
                  c.shortHash === props.commit || c.hash === props.commit
                }
              />
              <label htmlFor={c.hash}>
                <span className={styles.hash}>{c.shortHash}</span>
                <span className={styles.message}>{c.shortMessage}</span>
              </label>
            </div>
          ))
        ) : (
          <></>
        )}
      </fieldset>
    </>
  );
}
