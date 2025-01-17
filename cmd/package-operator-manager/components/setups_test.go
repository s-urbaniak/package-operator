package components

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	ctrl "sigs.k8s.io/controller-runtime"

	manifestsv1alpha1 "package-operator.run/apis/manifests/v1alpha1"
)

var errTest = errors.New("test")

func Test_setupAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cm := &controllerMock{}

		cm.On("SetupWithManager", mock.Anything).
			Return(nil)

		err := setupAll(nil, []controllerSetup{
			{name: "test", controller: cm},
		})
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		cm := &controllerMock{}

		cm.On("SetupWithManager", mock.Anything).
			Return(errTest)

		err := setupAll(nil, []controllerSetup{
			{name: "test", controller: cm},
		})
		require.EqualError(
			t, err, "unable to create controller for test: test")
	})
}

func TestAllControllers(t *testing.T) {
	var mocks []*controllerMock
	newMock := func() *controllerMock {
		m := &controllerMock{}
		m.On("SetupWithManager", mock.Anything).
			Return(nil)
		mocks = append(mocks, m)
		return m
	}
	var (
		os     = newMock()
		cos    = newMock()
		osp    = newMock()
		cosp   = newMock()
		od     = newMock()
		cod    = newMock()
		pkg    = newMock()
		cpkg   = newMock()
		otmpl  = newMock()
		cotmpl = newMock()
	)
	all := AllControllers{
		ObjectSet:        ObjectSetController{os},
		ClusterObjectSet: ClusterObjectSetController{cos},

		ObjectSetPhase:        ObjectSetPhaseController{osp},
		ClusterObjectSetPhase: ClusterObjectSetPhaseController{cosp},

		ObjectDeployment:        ObjectDeploymentController{od},
		ClusterObjectDeployment: ClusterObjectDeploymentController{cod},

		Package:        PackageController{pkg},
		ClusterPackage: ClusterPackageController{cpkg},

		ObjectTemplate:        ObjectTemplateController{otmpl},
		ClusterObjectTemplate: ClusterObjectTemplateController{cotmpl},
	}
	err := all.SetupWithManager(nil)
	require.NoError(t, err)

	for _, m := range mocks {
		m.AssertExpectations(t)
	}
	assert.Len(t, all.List(), 10)
}

func TestBootstrapControllers(t *testing.T) {
	var mocks []*controllerMock
	newMock := func() *controllerMock {
		m := &controllerMock{}
		m.On("SetupWithManager", mock.Anything).
			Return(nil)
		mocks = append(mocks, m)
		return m
	}
	var (
		cos  = newMock()
		cod  = newMock()
		cpkg = newMock()
	)
	all := BootstrapControllers{
		ClusterObjectSet:        ClusterObjectSetController{cos},
		ClusterObjectDeployment: ClusterObjectDeploymentController{cod},
		ClusterPackage:          ClusterPackageController{cpkg},
	}
	err := all.SetupWithManager(nil)
	require.NoError(t, err)

	for _, m := range mocks {
		m.AssertExpectations(t)
	}
	assert.Len(t, all.List(), 3)
}

type controllerMock struct {
	mock.Mock
}

func (m *controllerMock) SetupWithManager(mgr ctrl.Manager) error {
	args := m.Called(mgr)
	return args.Error(0)
}

func (m *controllerMock) SetEnvironment(env *manifestsv1alpha1.PackageEnvironment) {
	m.Called(env)
}
