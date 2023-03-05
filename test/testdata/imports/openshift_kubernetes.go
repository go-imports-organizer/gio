package imports

import (
	"fmt"
	_ "fmt"
	_ "io"
	_ "reflect"
	_ "sort"

	_ "github.com/openshift/api/build/v1"
	_ "github.com/openshift/client-go/build/clientset/versioned/typed/build/v1"
	_ "github.com/openshift/imagebuilder"
	_ "github.com/openshift/imagebuilder/dockerfile/command"
	_ "github.com/openshift/imagebuilder/dockerfile/parser"
	_ "github.com/openshift/library-go/pkg/git"
	_ "github.com/openshift/library-go/pkg/image/reference"
	_ "github.com/openshift/source-to-image/pkg/scm/git"
	_ "github.com/openshift/source-to-image/pkg/util"
	_ "k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/api/errors"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/apimachinery/pkg/util/sets"
	_ "k8s.io/apimachinery/pkg/util/uuid"
	_ "k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	fmt.Print("Hello, World!")
}
