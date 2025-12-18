package server

import "testing"

func TestServerBuilder_BuildValidation(t *testing.T) {
	_, err := NewServerBuilder().SetControlAddress("127.0.0.1:9090").Build()
	if err == nil {
		t.Fatalf("expected error for missing server address")
	}

	_, err = NewServerBuilder().SetAddress("0.0.0.0:4311").Build()
	if err == nil {
		t.Fatalf("expected error for not giving controll address")
	}

	_, err = NewServerBuilder().
		SetAddress("0.0.0.0:4311").
		SetControlAddress("0.0.0.0:9090").
		SetTLS(TLSConfig{Cert: "x", Key: ""}).
		Build()
	if err == nil {
		t.Fatalf("expected error for missing TLS key")
	}
}
