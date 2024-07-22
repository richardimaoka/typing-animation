"use client";

// !!!!
// Can only be used via Next.js dynamic import with ssr false option,
// due to "monaco-editor" module using browser-side `navigator` inside.
// !!!!
import { editor } from "monaco-editor";
import { EditorBare } from "./EditorBare";
import { useEditSequence } from "./hooks/useEditSequence";
import { useEditorInstance } from "./hooks/useEditorInstance";
import { useEditorTextUpdate } from "./hooks/useEditorTextUpdate";
import { useLanguageUpdate } from "./hooks/useLanguageUpdate";
import { ReactNode, useState } from "react";
import styles from "./EditorInnerOnlyDynamicallyImportable.module.css";

interface Props {
  editorText: string;
  language: string;
  editSequence?: {
    id: string;
    edits: editor.IIdentifiedSingleEditOperation[];
  };
}

// `default` export, for easier use with Next.js dynamic import
export default function EditorInnerOnlyDynamicallyImportable(props: Props) {
  /**
   * Monaco editor instance and readiness
   */
  const [editorInstance, onMount] = useEditorInstance();
  // isReady = true, if the initial rendering is finished
  const [isReady, setIsReady] = useState(false);

  /**
   * Basic editor text and its language
   */
  useEditorTextUpdate(editorInstance, props.editorText);
  useLanguageUpdate(editorInstance, props.language);

  /**
   * Edits
   */
  const { isEditCompleted } = useEditSequence(
    editorInstance,
    props.editSequence
  );

  return (
    // Needs the outer <div> for bounding box size retrieval.
    <div className={styles.component}>
      {/* The outer <div> has to be separate from the inner <div> because
       ** the inner div has `display: "none"` at the beginning which makes
       ** the bounding box zero-sized.
       */}
      <div
        // Until initial rendering is done (i.e.) onChange is at least called once,
        // delay the display of this display-control component. This is necessary
        // because otherwise the monaco editor moves the carousel unexpectedly by
        // (seemingly) calling scrollIntoView().
        //
        // By setting `display: "none"`, scrollIntoView() will not take effect and
        // the carousel does not move.
        style={{ display: isReady ? "block" : "none" }}
        className={styles.displayControl}
      >
        <EditorBare
          onMount={onMount}
          onChange={() => {
            //onChange is called after the initial rendering
            setIsReady(true);
          }}
        />
      </div>
    </div>
  );
}
