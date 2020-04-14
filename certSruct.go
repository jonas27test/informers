package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/deprecated/scheme"
	"k8s.io/client-go/rest"
)

const GroupName = "cert-manager.io"
const GroupVersion = "v1alpha3"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&CMCertificate{},
	)

	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

// DeepCopyObject returns a generically typed copy of an object
func (in *CMCertificate) DeepCopyObject() runtime.Object {
	out := CMCertificate{}
	in.DeepCopyInto(&out)

	return &out
}

func (in *CMCertificate) DeepCopyInto(out *CMCertificate) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = CMCertSpec{
		SecretName: in.Spec.SecretName,
	}
}

func (c *certClient) Get(name string, opts metav1.GetOptions) {
	result := CMCertificate{}
	err := c.restClient.
		Get().
		Namespace("inf").
		Resource("projects").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(nil).
		Into(&result)

	return &result, err
}

type certClient struct {
	restClient rest.Interface
	ns         string
}

type CMCertificateList struct {
}

type CMCertificate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec CMCertSpec `json:"spec"`
}

type CMCertSpec struct {
	// Secret names are always required.
	SecretName   string   `json:"secretName,required"`
	Duration     string   `json:"duration,omitempty"`
	RenewBefore  string   `json:"renewBefore,omitempty"`
	Organization []string `json:"organization,omitempty"`
	// The use of the common name field has been deprecated since 2000 and is
	// discouraged from being used.
	CommonName   string   `json:"commonName,omitempty"`
	IsCA         bool     `json:"isCA,omitempty"`
	KeySize      int64    `json:"keySize,omitempty"`
	KeyAlgorithm string   `json:"keyAlgorithm,omitempty"`
	KeyEncoding  string   `json:"KeyEncoding,omitempty"`
	Usages       []string `json:"usages,omitempty"`
	// At least one of a DNS Name, USI SAN, or IP address is required.
	DNSNames    []string `json:"dnsNames,required"`
	URISANs     []string `json:"uriSANs,omitempty"`
	IPAddresses []string `json:"ipAddresses,omitempty"`
	// Issuer references are always required. (Either Cluster or NS)
	IssuerRef CMIssuerRef `json:"issuerRef,required"`
}

type CMIssuerRef struct {
	Name string `json:"name,required"`
	// We can reference ClusterIssuers by changing the kind here.
	// The default value is Issuer (i.e. a locally namespaced Issuer)
	Kind string `json:"kind,omitempty"`
	// This is optional since cert-manager will default to this value however
	// if you are using an external issuer, change this to that issuer group.
	Group string `json:"group,omitempty"`
}
