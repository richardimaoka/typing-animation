package gitpkg

import (
	"testing"
)

func Test(t *testing.T) {
	orgname := "spf13"
	reponame := "cobra"
	filePath := "command.go"

	repo, err := openOrClone(orgname, reponame)
	if err != nil {
		t.Fatal(err)
	}

	commitHashes := []string{
		"e94f6d0dd9a5e5738dca6bce03c4b1207ffbc0ec", //Address golangci-lint deprecation warnings, enable some more linters (#2152)
		"8003b74a10ef0d0d84fe3c408d3939d86fdeb210", //Remove fully inactivated linters (#2148)
		"4fb0a66a3436bd34b03b858c729404e99cd3124f", //flags: clarify documentation that LocalFlags related function do not modify the state (#2064)
		"a73b9c391a9489d20f5ee1480e75d3b99fc8c7e2", //Fix help text for runnable plugin command
		"df547f5fc6ee86071f73c36c16afd885bb4e3f28", //Fix help text for plugins
		"3d8ac432bdad89db04ab0890754b2444d7b4e1cf", //Micro-optimizations (#1957)
		"890302a35f578311404a462b3cdd404f34db3720", //Support usage as plugin for tools like kubectl (#2018)
		"8b1eba47616566fc4d258a93da48d5d8741865f0", //Fix linter errors (#2052)
		"4cafa37bc4bb85633b4245aa118280fe5a9edcd5", //Allow running persistent run hooks of all parents (#2044)
		"95d8a1e45d7719c56dc017e075d3e6099deba85d", //Add notes to doc on preRun and postRun condition (#2041)
		"0c72800b8dba637092b57a955ecee75949e79a73", //Customizable error message prefix (#2023)
		"285460dca6152bb86994fa4a9659c24ca4060e2f", //command: temporarily disable G602 due to securego/gosec#1005 (#2022)
		"9e6b58afc70c60a6b3c8a0138fb25acc734d47e3", //update copyright year (#1927)
		"b4f979ae352828a0153281b590771ce9588d67f8", //completions: do not detect arguments with dash as 2nd char as flag (#1817)
		"bf11ab6321f2831be5390fdd71ab0bb2edbd9ed5", //fix: func name in doc strings (#1885)
		"6b0bd3076cfafd1c108264ed1e4aa0c0fe3f8537", //fix: don't remove flag value that matches subcommand name (#1781)
		"10cf7be9972ee5b5994cf760ecba642309ac8685", //Check for group presence after full initialization (#1839)
		"860791844ed3a2e544a9b9bbbcb14144a948ad20", //feat: make InitDefaultCompletionCmd public (#1467)
		"2169adb5749372c64cdd303864ae8a444da6350f", //Add groups for commands in help (#1003)
		"212ea4078323771dc49b6f25a41d84efbaac3a4c", //Include --help and --version flag in completion (#1813)
		"93d1913fb03362f97e95aeacc7d1541764cafc2f", //Add OnFinalize method (#1788)
		"7039e1fa214cfc1de404ed6540158c8fda64a758", //Add '--version' flag to Help output (#1707)
		"fce8d8aeb08dc6afe413cc0af67a7fbb3cffec4c", //Expose ValidateRequiredFlags and ValidateFlagGroups (#1760)
		"6d978a911e7fff69b1ca2a873dd91d78ebca44cf", //add missing license headers (#1809)
		"d689184a421607457a18131e0a2b602fec22e3b4", //Support for case-insensitive command names (#1802)
		"2a7647ff4661fd5bc54bdc022d349d7efe29674f", //Clarify SetContext documentation (#1748)
		"22b617914c8890ba20db7ceafcdc2ef4ca4817d3", //fix: show flags that shadow parent persistent flag in child help (#1776)
		"b9ca5949e2f58373e8e4c3823c213401f7d9d0e3", //use errors.Is() to check for errors (#1730)
		"ea94a3db55f84f026891709d82ebf25e17f89e0d", //undefined or nil Args default to ArbitraryArgs (#1612)
		"68b6b24f0c9926a779f7113d1040396a30fdedaf", //Add ability to mark flags as required or exclusive as a group (#1654)
		"3a1795bc253b686e565770630a78baee3e7ae7f4", //Fix Command.Context comment (#1639)
		"f848943afd7212766aadc19256cf9e7384980281", //Add Command.SetContext (#1551)
		"de187e874d1ca382320088f8f6d76333408e5c2e", //Fix flag completion (#1438)
		"3c8a19ecd384dc4229545bb174310a50c493f4ae", //fix RegisterFlagCompletionFunc concurrent map writes error (#1423)
		"b36196066e3b97b3cc87a352c81279af77028cc8", //Bash completion V2 with completion descriptions (#1146)
		"6d00909120c77b54b0c9974a4e20ffc540901b98", //Pass context to completion (#1265)
		"b312f0a8ef6f7ceb211af700083046da437db67c", //Create 'completion' command automatically (#1192)
		"652c755d3751109d878a9c8228c0bcb7aed0d4f5", //Use golangci-lint (#1044)
		"7df62f7668c7fbeafd6dbb45ad0294f4ea08c0ec", //fix typos (#1274)
		"40d34bca1bffe2f5e84b18d7fd94d5b3c02275a6", //Fix stderr printing functions (#894)
		"8cfa4b4acf63ab83334178f96ac694d3ea24e231", //Add documentation for Use (#1188)
		"04318720db1743b8488c86b2f7dca6d9663cb2f2", //Add completion for help command (#1136)
		"5155946348eed0f79a76f7743407c0c933e3b5f0", //Ignore required flags when DisableFlagParsing (#1095)
		"b84ef40338051c81c2916fe90c7e73e6e5583182", //Rename BashCompDirectives to ShellCompDirectives (#1082)
		"b80aeb17fc46362ff9cea51437a719322f8965ac", //Add support for custom completions in Go (#1035)
		"6607e6b8603f56adb027298ee6695e06ffb3a819", //Partial Revert of #922 (#1068)
		"95f2f73ed97e57387762620175ccdc277a8597a0", //Add short version flag -v when not otherwise set (#996)
		"3c2624538b7d0935103b37a9313661ffaad30d46", //Correct documentation for InOrStdin (#929)
		"0da06874266c88228b8f14615396a1f6bfc90ed7", //Add support for context.Context
		"0d9d2d46f3099574b6d524c3d8834dd592a80fb3", //Revert change so help is printed on stdout again (#1004)
		"993cc5372a05240dfd59e3ba952748b36b2cd117", //Adjustments per PR review feedback from @bogem
		"51f06c7dd1e73470976107fc6931b21143b83676", //Correct all complaints from golint
		"9334a46bd6b3887f3561d705440038ec93b7f62e", //Return an error in the case of unrunnable subcommand
	}

	for i := 0; i < len(commitHashes)-1; i++ {
		t.Run(commitHashes[i], func(t *testing.T) {
			currentHash := commitHashes[i]
			_, err := FileContentsInCommit(repo, currentHash, filePath)
			if err != nil {
				t.Fatal(err)
			}

			nextHash := commitHashes[i+1]
			_, err = FileContentsInCommit(repo, nextHash, filePath)
			if err != nil {
				t.Fatal(err)
			}

			// edits :=
			// for edits
			//
		})
	}
}
