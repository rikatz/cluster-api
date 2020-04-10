package ignition

import (
	ignTypes "github.com/flatcar-linux/ignition/config/v2_3/types"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/secret"
	"encoding/json"

)

const (
	standardJoinCommand            = "kubeadm join --config /tmp/kubeadm-join-config.yaml %s"
	retriableJoinScriptName        = "/usr/local/bin/kubeadm-bootstrap-script"
	retriableJoinScriptOwner       = "root"
	retriableJoinScriptPermissions = "0755"
)

// BaseUserData is shared across all the various types of files written to disk.
type BaseUserData struct {
	Header               string
	PreKubeadmCommands   []string
	PostKubeadmCommands  []string
	AdditionalFiles      []bootstrapv1.File
	WriteFiles           []bootstrapv1.File
	Users                []bootstrapv1.User
	NTP                  *bootstrapv1.NTP
	ControlPlane         bool
	UseExperimentalRetry bool
	KubeadmCommand       string
	KubeadmVerbosity     string
}

// ControlPlaneInput defines the context to generate a controlplane instance user data.
type ControlPlaneInput struct {
	BaseUserData
	secret.Certificates

	ClusterConfiguration string
	InitConfiguration    string
}

func NewInitControlPlane(input *ControlPlaneInput) ([]byte, error) {
	var ignitionConfig ignTypes.Config
	ignitionConfig.Ignition.Version = "2.3.0"
	
	// TODO: Move to a proper function
	for _, userItem := range input.Users {
		userObj := &ignTypes.PasswdUser{
			Name:              userItem.Name,
			SSHAuthorizedKeys: []ignTypes.SSHAuthorizedKey{"lalala"},
			Groups: []ignTypes.Group{"wheel"}

		}
		append(ignitionConfig.Passwd, userObj)
	}
	
	/*input.Header = cloudConfigHeader
	input.WriteFiles = input.Certificates.AsFiles()
	input.WriteFiles = append(input.WriteFiles, input.AdditionalFiles...)
	userData, err := generate("InitControlplane", controlPlaneCloudInit, input)
	if err != nil {
		return nil, err
	}*/

	userData, err := json.Marshal(ignitionConfig)
	if err != nil {
		return nil, err
	}
	return userData, nil
}
