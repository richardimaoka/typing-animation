"use client";

import styles from "./RepositoryNameField.module.css";
import { useRouter, usePathname, useSearchParams } from "next/navigation";

interface Props {
  reponame?: string;
}

export function RepositoryNameField(props: Props) {
  // https://github.com/vercel/next.js/discussions/47583
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const updateSearchParams = (newValue: string) => {
    // now you got a read/write object
    const current = new URLSearchParams(Array.from(searchParams.entries())); // -> has to use this form

    const paramName = "reponame";
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

  return (
    <>
      <label htmlFor="reponame" className={styles.label}>
        <div>GitHub repository</div>
      </label>
      <input
        id="reponame"
        className={styles.input}
        defaultValue={props.reponame}
        placeholder="{reponame}"
        onBlur={(e) => {
          updateSearchParams(e.target.value.trim());
        }}
        onKeyDown={(e) => {
          if (e.key === "Enter" && !e.nativeEvent.isComposing) {
            updateSearchParams(e.currentTarget.value.trim());
          }
        }}
      />
    </>
  );
}
