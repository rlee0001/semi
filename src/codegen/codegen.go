package codegen

import (
	"io/ioutil"
	"os/exec"

	"semi/src/ast"
)

// TODO: Take a "globals" map? Or is that generated internally by codeGenWalker? Or should each node have its own ns?

func ToLlvmFile(node ast.Node) (string, error) {
	codeGenWalker := NewCodeGenWalker()

	_, err := codeGenWalker.VisitNode(node, nil)
	if err != nil {
		return "", err
	}

	llFile, err := ioutil.TempFile("", "*.ll")
	if err != nil {
		return "", err
	}

	_, err = codeGenWalker.getIrModule().WriteTo(llFile)
	if err != nil {
		return "", err
	}

	if err := llFile.Close(); err != nil {
		return "", err
	}

	return llFile.Name(), nil
}

func RunClang(llvmFilePath string, binaryFileName string) error {
	_, err := exec.LookPath("clang")
	if err != nil {
		return err
	}

	err = exec.Command("clang", llvmFilePath, "-O3", "-o", binaryFileName).Run()
	if err != nil {
		return err
	}

	return nil
}
