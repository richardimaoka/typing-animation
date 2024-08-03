"use client";

import styles from "./OrgNameField.module.css";
import { useRouter, usePathname, useSearchParams } from "next/navigation";

interface Props {
  orgname: string | undefined;
}

export function OrgNameField(props: Props) {
  // https://github.com/vercel/next.js/discussions/47583
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const updateSearchParams = (newValue: string) => {
    // now you got a read/write object
    const current = new URLSearchParams(Array.from(searchParams.entries())); // -> has to use this form

    // update as necessary
    if (!newValue) {
      current.delete("selected");
    } else {
      current.set("selected", newValue);
    }

    // cast to string
    const search = current.toString();
    // or const query = `${'?'.repeat(search.length && 1)}${search}`;
    const query = search ? `?${search}` : "";

    router.push(`${pathname}${query}`);
  };

  return (
    <>
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
