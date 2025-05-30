package put_branch_in_git

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	codebaseApi "github.com/epam/edp-codebase-operator/v2/api/v1"
	"github.com/epam/edp-codebase-operator/v2/controllers/codebasebranch/chain"
	"github.com/epam/edp-codebase-operator/v2/controllers/codebasebranch/service"
	gitServerMocks "github.com/epam/edp-codebase-operator/v2/pkg/git/mocks"
	"github.com/epam/edp-codebase-operator/v2/pkg/util"
)

const (
	fakeName      = "fake-name"
	fakeNamespace = "fake-namespace"
)

func TestPutBranchInGit_ShouldBeExecutedSuccessfullyWithDefaultVersioning(t *testing.T) {
	c := &codebaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "default",
		},
		Spec: codebaseApi.CodebaseSpec{
			GitServer:  "gitserver",
			GitUrlPath: "/test-app",
		},
		Status: codebaseApi.CodebaseStatus{
			Available: true,
		},
	}

	gitUser := "git-user"
	gs := &codebaseApi.GitServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitserver",
			Namespace: "default",
		},
		Spec: codebaseApi.GitServerSpec{
			NameSshKeySecret: "secret",
			GitHost:          "git-host",
			SshPort:          22,
			GitUser:          gitUser,
		},
	}

	sshKey := "fake"
	s := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			util.PrivateSShKeyName: []byte(sshKey),
		},
	}

	cb := &codebaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "feature-branch",
			Namespace: "default",
		},
		Spec: codebaseApi.CodebaseBranchSpec{
			CodebaseName: "test-app",
			BranchName:   "feature-branch",
			FromCommit:   "commitsha",
		},
	}

	scheme := runtime.NewScheme()
	require.NoError(t, codebaseApi.AddToScheme(scheme))
	require.NoError(t, coreV1.AddToScheme(scheme))

	fakeCl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(c, gs, cb, s).
		WithStatusSubresource(cb).
		Build()

	mGit := gitServerMocks.NewMockGit(t)

	mGit.On(
		"CloneRepositoryBySsh",
		testifymock.Anything,
		sshKey,
		gs.Spec.GitUser,
		testifymock.Anything,
		testifymock.Anything,
		gs.Spec.SshPort,
	).Return(nil)
	mGit.On(
		"GetCurrentBranchName",
		testifymock.Anything,
	).Return("default-branch", nil)
	mGit.On(
		"CreateRemoteBranch",
		sshKey,
		gs.Spec.GitUser,
		testifymock.Anything,
		cb.Spec.BranchName,
		cb.Spec.FromCommit,
		gs.Spec.SshPort,
	).Return(nil)
	mGit.On(
		"CheckoutRemoteBranchBySSH",
		sshKey,
		gs.Spec.GitUser,
		testifymock.Anything,
		c.Spec.DefaultBranch,
	).Return(nil)

	err := PutBranchInGit{
		Client: fakeCl,
		Git:    mGit,
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), cb)

	assert.NoError(t, err)
}

func TestPutBranchInGit_ShouldFailgetCurrentbranch(t *testing.T) {
	c := &codebaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-app",
			Namespace: "default",
		},
		Spec: codebaseApi.CodebaseSpec{
			GitServer:  "gitserver",
			GitUrlPath: "/test-app",
		},
		Status: codebaseApi.CodebaseStatus{
			Available: true,
		},
	}

	gitUser := "git-user"
	gs := &codebaseApi.GitServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "gitserver",
			Namespace: "default",
		},
		Spec: codebaseApi.GitServerSpec{
			NameSshKeySecret: "secret",
			GitHost:          "git-host",
			SshPort:          22,
			GitUser:          gitUser,
		},
	}

	sshKey := "fake"
	s := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "default",
		},
		Data: map[string][]byte{
			util.PrivateSShKeyName: []byte(sshKey),
		},
	}

	cb := &codebaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "feature-branch",
			Namespace: "default",
		},
		Spec: codebaseApi.CodebaseBranchSpec{
			CodebaseName: "test-app",
			BranchName:   "feature-branch",
			FromCommit:   "commitsha",
		},
	}

	scheme := runtime.NewScheme()
	require.NoError(t, codebaseApi.AddToScheme(scheme))
	require.NoError(t, coreV1.AddToScheme(scheme))

	fakeCl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(c, gs, cb, s).
		WithStatusSubresource(cb).
		Build()

	mGit := gitServerMocks.NewMockGit(t)

	mGit.On(
		"CloneRepositoryBySsh",
		testifymock.Anything,
		sshKey,
		gs.Spec.GitUser,
		testifymock.Anything,
		testifymock.Anything,
		gs.Spec.SshPort,
	).Return(nil)
	mGit.On(
		"GetCurrentBranchName",
		testifymock.Anything,
	).Return("", errors.New("failed to get current branch"))

	err := PutBranchInGit{
		Client: fakeCl,
		Git:    mGit,
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), cb)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get current branch")
}

func TestPutBranchInGit_ShouldFailCreateRemoteBranch(t *testing.T) {
	t.Setenv(util.WorkDirEnv, t.TempDir())

	c := &codebaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseSpec{
			GitServer:     fakeName,
			GitUrlPath:    fakeName,
			DefaultBranch: "main",
		},
		Status: codebaseApi.CodebaseStatus{
			Available: true,
		},
	}

	gs := &codebaseApi.GitServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.GitServerSpec{
			NameSshKeySecret: fakeName,
			GitHost:          fakeName,
			SshPort:          22,
			GitUser:          fakeName,
		},
	}

	s := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Data: map[string][]byte{
			"keyName": []byte("fake"),
		},
	}

	cb := &codebaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseBranchSpec{
			CodebaseName: fakeName,
			BranchName:   fakeName,
			FromCommit:   "commitsha",
		},
	}

	scheme := runtime.NewScheme()
	require.NoError(t, codebaseApi.AddToScheme(scheme))
	require.NoError(t, coreV1.AddToScheme(scheme))

	fakeCl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(c, gs, cb, s).
		WithStatusSubresource(cb).
		Build()

	mGit := gitServerMocks.NewMockGit(t)

	mGit.On(
		"CloneRepositoryBySsh",
		testifymock.Anything,
		testifymock.Anything,
		testifymock.Anything,
		testifymock.Anything,
		testifymock.Anything,
		testifymock.Anything,
	).Return(nil)

	mGit.On(
		"GetCurrentBranchName",
		testifymock.Anything,
	).Return("main", nil)

	mGit.On(
		"CreateRemoteBranch",
		testifymock.Anything,
		testifymock.Anything,
		testifymock.Anything,
		fakeName,
		testifymock.Anything,
		testifymock.Anything,
	).Return(errors.New("failed to create remote branch"))

	err := PutBranchInGit{
		Client: fakeCl,
		Git:    mGit,
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), cb)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create remote branch")
}

func TestPutBranchInGit_CodebaseShouldNotBeFound(t *testing.T) {
	cb := &codebaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseBranchSpec{
			CodebaseName: fakeName,
		},
	}
	scheme := runtime.NewScheme()
	require.NoError(t, codebaseApi.AddToScheme(scheme))

	fakeCl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(cb).
		WithStatusSubresource(cb).
		Build()

	err := PutBranchInGit{
		Client: fakeCl,
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), cb)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch Codebase")
	assert.Equal(t, codebaseApi.PutGitBranch, cb.Status.Action)
}

func TestPutBranchInGit_ShouldThrowCodebaseBranchReconcileError(t *testing.T) {
	c := &codebaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseSpec{
			GitServer:  fakeName,
			GitUrlPath: fakeName,
		},
		Status: codebaseApi.CodebaseStatus{
			Available: false,
		},
	}

	cb := &codebaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseBranchSpec{
			CodebaseName: fakeName,
		},
	}
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(metav1.SchemeGroupVersion, c, cb)
	fakeCl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(c, cb).
		WithStatusSubresource(cb).
		Build()

	err := PutBranchInGit{
		Client: fakeCl,
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), cb)

	_, ok := err.(*util.CodebaseBranchReconcileError)
	assert.True(t, ok, "wrong type of error")
}

func TestPutBranchInGit_ShouldBeExecutedSuccessfullyWithEdpVersioning(t *testing.T) {
	c := &codebaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseSpec{
			GitServer:  fakeName,
			GitUrlPath: fakeName,
			Versioning: codebaseApi.Versioning{
				Type:      codebaseApi.VersioningTypeSemver,
				StartFrom: nil,
			},
			DefaultBranch: "main",
		},
		Status: codebaseApi.CodebaseStatus{
			Available: true,
		},
	}

	gs := &codebaseApi.GitServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.GitServerSpec{
			NameSshKeySecret: fakeName,
			GitHost:          fakeName,
			SshPort:          22,
			GitUser:          fakeName,
		},
	}

	s := &coreV1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Data: map[string][]byte{
			"keyName": []byte("fake"),
		},
	}

	cb := &codebaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseBranchSpec{
			CodebaseName: fakeName,
			Version:      util.GetStringP("version3"),
			BranchName:   fakeName,
			FromCommit:   "",
		},
		Status: codebaseApi.CodebaseBranchStatus{
			VersionHistory: []string{"version1", "version2"},
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(metav1.SchemeGroupVersion, c, gs, cb)
	scheme.AddKnownTypes(coreV1.SchemeGroupVersion, s)
	fakeCl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(c, gs, cb, s).
		WithStatusSubresource(cb).
		Build()

	mGit := gitServerMocks.NewMockGit(t)

	port := int32(22)
	wd := chain.GetCodebaseBranchWorkingDirectory(cb)

	repoSshUrl := util.GetSSHUrl(gs, c.Spec.GetProjectID())
	mGit.On("CloneRepositoryBySsh", testifymock.Anything, "", fakeName, repoSshUrl, wd, port).
		Return(nil)
	mGit.On(
		"GetCurrentBranchName",
		testifymock.Anything,
	).Return("main", nil)
	mGit.On("CreateRemoteBranch", "", fakeName, wd, fakeName, "", port).Return(nil)

	err := PutBranchInGit{
		Client: fakeCl,
		Git:    mGit,
		Service: &service.CodebaseBranchServiceProvider{
			Client: fakeCl,
		},
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), cb)

	assert.NoError(t, err)
}

func TestPutBranchInGit_ShouldFailToSetIntermediateStatus(t *testing.T) {
	cb := &codebaseApi.CodebaseBranch{}

	scheme := runtime.NewScheme()
	fakeCl := fake.NewClientBuilder().WithScheme(scheme).WithRuntimeObjects().Build()

	err := PutBranchInGit{
		Client: fakeCl,
		Service: &service.CodebaseBranchServiceProvider{
			Client: fakeCl,
		},
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), cb)

	assert.Error(t, err)
}

func TestPutBranchInGit_GitServerShouldNotBeFound(t *testing.T) {
	c := &codebaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseSpec{
			GitServer:  fakeName,
			GitUrlPath: fakeName,
		},
		Status: codebaseApi.CodebaseStatus{
			Available: true,
		},
	}

	cb := &codebaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseBranchSpec{
			CodebaseName: fakeName,
		},
	}

	scheme := runtime.NewScheme()
	require.NoError(t, codebaseApi.AddToScheme(scheme))
	require.NoError(t, coreV1.AddToScheme(scheme))
	fakeCl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(c, cb).
		WithStatusSubresource(cb).
		Build()

	err := PutBranchInGit{
		Client: fakeCl,
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), cb)

	assert.Error(t, err)

	assert.Contains(t, err.Error(), "failed to fetch GitServer")
}

func TestPutBranchInGit_SecretShouldNotBeFound(t *testing.T) {
	c := &codebaseApi.Codebase{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseSpec{
			GitServer:  fakeName,
			GitUrlPath: fakeName,
		},
		Status: codebaseApi.CodebaseStatus{
			Available: true,
		},
	}

	gs := &codebaseApi.GitServer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.GitServerSpec{
			NameSshKeySecret: fakeName,
			GitHost:          fakeName,
			SshPort:          22,
			GitUser:          fakeName,
		},
	}

	cb := &codebaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeName,
			Namespace: fakeNamespace,
		},
		Spec: codebaseApi.CodebaseBranchSpec{
			CodebaseName: fakeName,
		},
	}

	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(metav1.SchemeGroupVersion, c, gs, cb)
	scheme.AddKnownTypes(coreV1.SchemeGroupVersion, &coreV1.Secret{})
	fakeCl := fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(c, gs, cb).
		WithStatusSubresource(cb).
		Build()

	err := PutBranchInGit{
		Client: fakeCl,
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), cb)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), fmt.Sprintf("failed to get %s secret", fakeName))
}

func TestPutBranchInGit_SkipAlreadyCreated(t *testing.T) {
	codeBaseBranch := &codebaseApi.CodebaseBranch{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "fake",
			Namespace: "default",
		},
		Spec: codebaseApi.CodebaseBranchSpec{
			CodebaseName: "fake",
			BranchName:   "fake",
		},
		Status: codebaseApi.CodebaseBranchStatus{
			Git: codebaseApi.CodebaseBranchGitStatusBranchCreated,
		},
	}

	err := PutBranchInGit{
		Client: fake.NewClientBuilder().Build(),
		Git:    gitServerMocks.NewMockGit(t),
	}.ServeRequest(ctrl.LoggerInto(context.Background(), logr.Discard()), codeBaseBranch)

	require.NoError(t, err)
}
