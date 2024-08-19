"use client";

import styles from "./FilePathSelection.module.css";
import { useRouter, usePathname, useSearchParams } from "next/navigation";

interface Props {
  filepath?: string;
  fileSelection?: string[];
}

export function FilePathSelection(props: Props) {
  // https://github.com/vercel/next.js/discussions/47583
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const updateSearchParams = (newValue: string) => {
    // now you got a read/write object
    const current = new URLSearchParams(Array.from(searchParams.entries())); // -> has to use this form

    // update as necessary
    const paramName = "filepath";
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

  const defaultValue =
    props.filepath && props.fileSelection?.includes(props.filepath)
      ? props.filepath
      : "unselected";

  return (
    <>
      <label htmlFor="filepath" className={styles.label}>
        file path
      </label>
      <select
        id="filepath"
        onChange={(e) => {
          updateSearchParams(e.target.value.trim());
        }}
        defaultValue={defaultValue}
      >
        {props.fileSelection && props.fileSelection.length > 0 ? (
          <>
            <option value="unselected" disabled>
              choose an option below
            </option>
            {props.fileSelection.map((filepath) => (
              <option key={filepath} value={filepath}>
                {filepath}
              </option>
            ))}
          </>
        ) : (
          <option disabled>no option is available</option>
        )}
      </select>
    </>
  );
}
