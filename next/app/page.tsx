import { Left } from "./components/Left";
import styles from "./page.module.css";
import { promises as fs } from "fs";

interface Props {
  params: { slug: string }; // for dynamic routes only
  searchParams: {
    orgname?: string | string[];
    reponame?: string | string[];
    branch?: string | string[];
    filepath?: string | string[];
  };
}

function retrieveParam(
  param: string | string[] | undefined
): string | undefined {
  return typeof param === "string" && param !== "" ? param : undefined;
}

export default async function Page(props: Props) {
  const orgname = retrieveParam(props.searchParams.orgname);
  const reponame = retrieveParam(props.searchParams.reponame);
  const branch = retrieveParam(props.searchParams.branch) || "main";
  const filepath = retrieveParam(props.searchParams.filepath) || "main";

  const dataFile = "/app/data/files.json";
  const fileContents = await fs
    .readFile(process.cwd() + dataFile, "utf8")
    .catch((x) => {
      console.log(x);
      return undefined;
    });

  let files: string[] | undefined = undefined;
  if (fileContents) {
    try {
      files = JSON.parse(fileContents) as string[];
    } catch (error) {
      console.log(`JSON parse error on file = ${dataFile}`, error);
      files = undefined;
    }
  }

  return (
    <div className={styles.component}>
      <Left
        orgname={orgname}
        reponame={reponame}
        branch={branch}
        files={files}
        filepath={filepath}
        commits={[
          { hash: "sad65sd234f", message: "this is the first commit" },
          { hash: "asd897sdf87", message: "this is the second commit" },
          { hash: "890908sdr49", message: "this is the third commit" },
          { hash: "78searoa36a", message: "this is the fourth commit" },
          { hash: "905sfdjo439", message: "this is the fifth commit" },
        ]}
      />
      <div>right</div>
    </div>
  );
}
