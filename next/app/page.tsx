import { Left } from "./components/Left";
import styles from "./page.module.css";

export default function Page() {
  return (
    <div className={styles.component}>
      <Left />
      <div>right</div>
    </div>
  );
}
