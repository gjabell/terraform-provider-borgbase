package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SshKeyModel struct {
	AddedAt    types.String `tfsdk:"added_at"`
	Bits       types.Int64  `tfsdk:"bits"`
	HashMd5    types.String `tfsdk:"hash_md5"`
	HashSha256 types.String `tfsdk:"hash_sha256"`
	Id         types.String `tfsdk:"id"`
	LastUsedAt types.String `tfsdk:"last_used_at"`
	Name       types.String `tfsdk:"name"`
	PublicKey  types.String `tfsdk:"public_key"`
	Type       types.String `tfsdk:"type"`
}

func (m *SshKeyModel) update(key SshKeyPayload) {
	m.AddedAt = types.StringValue(key.AddedAt)
	m.Bits = types.Int64Value(int64(key.Bits))
	m.HashMd5 = types.StringValue(key.HashMd5)
	m.HashSha256 = types.StringValue(key.HashSha256)
	m.Id = types.StringValue(key.Id)
	m.LastUsedAt = types.StringValue(key.LastUsedAt)
	m.Name = types.StringValue(key.Name)

	publicKey := key.KeyData
	if key.Comment != "" {
		publicKey += " " + key.Comment
	}
	m.PublicKey = types.StringValue(publicKey)

	m.Type = types.StringValue(key.KeyType)
}

type SshKeyPayload struct {
	AddedAt    string `json:"addedAt"`
	Bits       int    `json:"bits"`
	HashMd5    string `json:"hashMd5"`
	HashSha256 string `json:"hashSha256"`
	Id         string `json:"id"`
	LastUsedAt string `json:"lastUsedAt"`
	Name       string `json:"name"`
	KeyData    string `json:"keyData"`
	KeyType    string `json:"keyType"`
	Comment    string `json:"comment"`
}

type SshKeysPayload []SshKeyPayload

type SshAddPayload struct {
	KeyAdded SshKeyPayload `json:"keyAdded"`
}

type SshDeletePayload struct {
	Ok bool `json:"ok"`
}
