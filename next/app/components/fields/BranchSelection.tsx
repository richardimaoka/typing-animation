"use client";

import styles from "./BranchSelection.module.css";
import { useRouter, usePathname, useSearchParams } from "next/navigation";

interface Props {
  branch?: string;
  brancheSelection?: string[];
}

export function BranchSelection(props: Props) {
  // https://github.com/vercel/next.js/discussions/47583
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const updateSearchParams = (newValue: string) => {
    // now you got a read/write object
    const current = new URLSearchParams(Array.from(searchParams.entries())); // -> has to use this form

    // update as necessary
    const paramName = "branch";
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
    props.branch && props.brancheSelection?.includes(props.branch)
      ? props.branch
      : "unselected";

  return (
    <>
      <label htmlFor="branch" className={styles.label}>
        branch
      </label>
      <select
        id="branch"
        onChange={(e) => {
          console.log("branch", e.target.value.trim());
          updateSearchParams(e.target.value.trim());
        }}
        defaultValue={defaultValue}
      >
        {props.brancheSelection && props.brancheSelection.length > 0 ? (
          <>
            <option value="unselected" disabled>
              choose an option below
            </option>
            {props.brancheSelection.map((branchName) => (
              <option key={branchName} value={branchName}>
                {branchName}
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
